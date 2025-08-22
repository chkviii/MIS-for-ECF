package handler

import (
	"net/http"
	"strconv"

	"crypto/ed25519"

	"mypage-backend/internal/config"
	"mypage-backend/internal/repo"
	"mypage-backend/internal/service"
)

// login handler
func LoginHandlerFunc(w http.ResponseWriter, r *http.Request) {
	//get login type username and password from the form
	loginType := r.FormValue("login_type")
	salt := r.FormValue("salt")
	sig := r.FormValue("password")
	username := r.FormValue(loginType)

	// Initialize database connection
	db, err := repo.InitDB(config.GlobalConfig)
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	userRepo := repo.NewUserRepository(db)

	var userID uint64
	if loginType == "uid" {
		// Convert username to uint for user ID
		userID, err = strconv.ParseUint(username, 10, 32)
		if err != nil {
			http.Error(w, "Invalid user ID format", http.StatusBadRequest)
			return
		}
	}

	var user *repo.User
	// Check if the user exists based on the login type
	switch loginType {
	case "username":
		user, err = userRepo.GetByUsername(username)
	case "email":
		user, err = userRepo.GetByEmail(username)
	case "uid":
		user, err = userRepo.GetByID(uint(userID))
	default:
		http.Error(w, "Invalid login type", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if ed25519.Verify([]byte(user.Password), []byte(salt), []byte(sig)) {
		s := service.NewSession(service.Sess, string(user.ID))
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
