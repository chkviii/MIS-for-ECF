package handler

import (
	"net/http"
	
	"mypage-backend/internal/util"
)

func HomeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the index.html file
		http.ServeFile(w, r, util.Html_Path("index.html"))
	})
}