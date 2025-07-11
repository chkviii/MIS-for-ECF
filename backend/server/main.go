package main

import (
	"fmt"
	"net/http"

	"mypage-backend/internal/config"
	"mypage-backend/internal/handler"
	// "mypage-backend/internal/middleware"
	// "mypage-backend/internal/repo"
	// "mypage-backend/internal/service"

	"github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {

	// Load configuration
	cfg := config.Load()

	fmt.Println("GO: Loaded configuration:", cfg)
	
	//main router
	r := chi.NewRouter()

	//main router config
	r.Use(middleware.Logger) // log the start and end of each request

	// Static file serving
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.Static_Path)))
	r.Handle("/static/*", fs)
	
	//routers
	r.Handle("/", handler.HomeHandler()) // Home page handler
	r.Handle("/blog", handler.BlogHandler()) // Blog page 
	r.Handle("/blog/{postID}", handler.BlogHandler()) // Blog post page with postID /blog/123?id=123
	r.Handle("/api/v0/post/{postID}", handler.BlogPostHandler()) // Blog post handler

	//listen and serve on port 33031
	fmt.Println("GO: Server starting on port", "http://localhost:33031")
	http.ListenAndServe(":33031", r)
	
}