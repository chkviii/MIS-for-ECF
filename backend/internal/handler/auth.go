package handler

import (
	"net/http"

	"mypage-backend/internal/repo"
	"mypage-backend/internal/service"
	"mypage-backend/internal/util"
)

func PreloginHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var req service.PreLoginRequest
	err := util.ParseJSONBody(r, &req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// handle prelogin logic
	userSvc := service.NewUserService(repo.GetUserRepo(), "")
	resp, err := userSvc.PreLogin(req)
	if err != nil {
		http.Error(w, "Prelogin failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// save session id  and expiration time in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    resp.SessID,
		HttpOnly: true,
		Secure:   true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "session_expire",
		Value:    resp.SessExpire,
		HttpOnly: true,
		Secure:   true,
	})

	// send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"session_id":"` + resp.SessID + `","server_pubkey":"` + resp.SrvPubKey + `","challenge":"` + resp.Challenge + `"}`))
}

// register page init
func RegisterHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the register.html file
		http.ServeFile(w, r, util.Html_Path("register.html"))
		// send server public key to the client
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"pubkey":"` + string(util.SrvEncPubKey()) + `"}`))

	})
}

// register handler
func RegisterHandlerFunc(w http.ResponseWriter, r *http.Request) {
}

// login handler
func LoginHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	err := util.ParseJSONBody(r, &req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// get session id from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Session cookie not found", http.StatusUnauthorized)
		return
	}
	sessID := cookie.Value

	// handle login logic
	userSvc := service.NewUserService(repo.GetUserRepo(), "")
	resp, err := userSvc.Login(sessID, req)

	if err != nil {
		http.Error(w, "Login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if !resp.Success {
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	} else {
		// update session cookie expiration time
		http.SetCookie(w, &http.Cookie{
			Name:     "session_expire",
			Value:    resp.SessExpire,
			HttpOnly: true,
			Secure:   true,
		})

		// send response
		w.WriteHeader(http.StatusOK)
	}
}

// logout handler
func LogoutHandlerFunc(w http.ResponseWriter, r *http.Request) {
}
