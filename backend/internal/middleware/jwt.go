package midware

import (
	"net/http"
)

// AuthMiddleware is a middleware that checks if the request has a valid session key.

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//
		sessionKey := r.Header.Get("Session-Key")
		if sessionKey == "" {
			http.Error(w, "Unauthorized: No session key provided", http.StatusUnauthorized)
			return
		}

		// Validate the session key (this is a placeholder, implement your own logic)
		if !isValidSessionKey(sessionKey) {
			http.Error(w, "Unauthorized: Invalid session key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isValidSessionKey(sessionKey string) bool {
	// Placeholder for session key validation logic
	// You can implement your own logic to check if the session key is valid
	// For example, you can check against a database or an in-memory store

	// _, err := service.Sess.Load(sessionKey)
	// return err == nil
	return true // For now, assume all session keys are valid
}
