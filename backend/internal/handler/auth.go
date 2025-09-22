package handler

import (
	"net/http"
	"strings"

	"crypto/ed25519"

	"mypage-backend/internal/repo"
	"mypage-backend/internal/service"
	"mypage-backend/internal/util"
)

// register page init
func RegisterHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the register.html file
		http.ServeFile(w, r, util.Html_Path("register.html"))
		// send server public key to the client
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"pubkey":"` + string(util.SrvPubKey()) + `"}`))

	})
}

// register handler
func RegisterHandlerFunc(w http.ResponseWriter, r *http.Request) {
}

// login handler
func LoginHandlerFunc(w http.ResponseWriter, r *http.Request) {
	//get login type username and password from the form
	loginType := r.FormValue("login_type")
	salt := r.FormValue("salt")
	sig := r.FormValue("password")
	username := r.FormValue(loginType)

	// Initialize database connection
	db, err := repo.GetDB()
	if err != nil {
		http.Error(w, "It's not your fault. \n Database Offline", http.StatusInternalServerError)
		return
	}

	userRepo := repo.GetUserRepo(db)

	var user *repo.User
	// Check if the user exists based on the login type
	switch loginType {
	case "username":
		user, err = userRepo.GetByUsername(username)
	case "email":
		user, err = userRepo.GetByEmail(strings.ToLower(username))
	case "uid":
		user, err = userRepo.GetByID(username)
	default:
		http.Error(w, "Invalid login type", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if ed25519.Verify([]byte(user.Password), []byte(salt), []byte(sig)) {
		s := service.SessMgr.NewSession(string(user.ID))
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    s.ID,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "session_expire",
			Value:    string(s.ExpiresAt),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		})
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

}

// logout handler
func LogoutHandlerFunc(w http.ResponseWriter, r *http.Request) {
}
