package router

import (
	"main/internal/middlewares"
	"main/internal/services"
	"net/http"
)

func userRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.Handle("GET /user/me", middlewares.Auth(http.HandlerFunc(services.UserProfile)))
	mux.HandleFunc("GET /user/blogs", services.AllBlogs)
	mux.Handle("POST /user/blogs", middlewares.Auth(http.HandlerFunc(services.CreateBlog)))
	mux.HandleFunc("GET /user/blogs/{blogId}", services.BlogByID)
	mux.Handle("PUT /user/blogs/{blogId}", middlewares.Auth(http.HandlerFunc(services.UpdateBlog)))
	mux.Handle("DELETE /user/blogs/{blogId}", middlewares.Auth(http.HandlerFunc(services.DeletBlogByID)))

	mux.HandleFunc("GET /user/blogs/{blogId}/comments", services.GetAllBlogComments)
	mux.Handle("POST /user/blogs/{blogId}/comments", middlewares.Auth(http.HandlerFunc(services.AddBlogComment)))
	mux.Handle("PUT /user/blogs/{blogId}/comments/{commentId}", middlewares.Auth(http.HandlerFunc(services.UpdateBlogComment)))
	mux.Handle("DELETE /user/blogs/{blogId}/comments/{commentId}", middlewares.Auth(http.HandlerFunc(services.DeleteBlogComment)))

	return mux
}
