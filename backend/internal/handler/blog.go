package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"mypage-backend/internal/util"
)

// blog handler serves the blog page
func BlogHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check if postID in the context
		postID := chi.URLParam(r, "postID")
		if postID != "" && true /* TODO: check if exist*/ {

			// Serve the blog.html file
			http.ServeFile(w, r, util.Html_Path("blog.html"))
		} else {
			//redirect to home page if no postID
			// context := context.WithValue(r.Context(), "error", "No postID provided")
			http.Error(w, "No postID provided", http.StatusBadRequest)
			// http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})
}

//api for blog post
func BlogPostHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the blog post file based on the URL parameter
		postID := chi.URLParam(r, "postID")
		
		//get article metadata

		//just write a simple response for now
		articleData := struct {
			PostID    string
			Title     string
			AuthorID  string
			AuthorName string
			LastUpdated string
			ContentFilePath string
		}{
			PostID:    postID,
			Title:     "示例文章标题",
			AuthorID:  "author123",
			AuthorName: "作者姓名",
			LastUpdated: "2023-10-01",
			ContentFilePath: "/static/article/" + postID + ".md",
		}

		//return a json response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(articleData); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}