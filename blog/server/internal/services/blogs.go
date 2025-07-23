package services

import (
	"encoding/json"
	"fmt"
	"main/internal/contexter"
	"main/internal/database"
	"main/pkg/response"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type Blog struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func AllBlogs(w http.ResponseWriter, r *http.Request) {
	allBlogsQuery := "SELECT id, title, body FROM blogs"
	rows, err := database.DB.Query(allBlogsQuery)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()
	var blogList []Blog

	for rows.Next() {
		var blog Blog
		err := rows.Scan(&blog.ID, &blog.Title, &blog.Body)
		if err != nil {
			response.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		blogList = append(blogList, blog)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]Blog{
		"data": blogList,
	})
}

func BlogByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	blogByIDQuery := "SELECT id, title, body FROM blogs WHERE id = $1"
	var blog Blog
	err := database.DB.QueryRow(blogByIDQuery, id).Scan(&blog.ID, &blog.Title, &blog.Body)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data": map[string]Blog{
			"blog": blog,
		},
	})
}

func CreateBlog(w http.ResponseWriter, r *http.Request) {
	user, ok := contexter.UserFromContext(r.Context())
	if !ok {
		response.Error(w, "please authenticate to create a post", http.StatusUnauthorized)
		return
	}

	type Blog struct {
		Title string `json:"title" validate:"required,min=4,max=250"`
		Body  string `json:"body" validate:"required,min=4,max=1000"`
	}

	var blog Blog
	json.NewDecoder(r.Body).Decode(&blog)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(&blog)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createBlogQuery := `
	INSERT INTO blogs (author_id, title, body, create_at, updated_at)
	VALUES( ( SELECT id FROM users WHERE email = $1), $2, $3, $4, $5)
	`
	date := time.Now().UTC()
	rows, err := database.DB.Query(createBlogQuery, user.Email, blog.Title, blog.Body, date, date)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer rows.Close()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Blog Successfully created",
	})
}

func UpdateBlog(w http.ResponseWriter, r *http.Request) {
	user, ok := contexter.UserFromContext(r.Context())
	if !ok {
		response.Error(w, "please authenticate to create a post", http.StatusUnauthorized)
		return
	}

	type Blog struct {
		Title string `json:"title" validate:"min=4,max=250"`
		Body  string `json:"body" validate:"min=4,max=1000"`
	}
	var blog Blog
	json.NewDecoder(r.Body).Decode(&blog)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(&blog)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if blog.Title == "" && blog.Body == "" {
		response.Error(w, "either title or body must be provided", http.StatusBadRequest)
		return
	}

	setQuery := "SET"

	if blog.Title != "" {
		setQuery += fmt.Sprintf("title = %s", blog.Title)
	}

	if blog.Body != "" {
		setQuery += fmt.Sprintf("body = %s", blog.Body)
	}

	updateQuery := "UPDATE blogs" + " " + setQuery + " " + fmt.Sprintf("updated_at = %s", time.Now().UTC().String()) + " WHERE author_id = ( SELECT * FROM users WHERE email = $1)"
	_, err = database.DB.Query(updateQuery, user.Email)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Blog Updated Successfully",
	})
}

func DeletBlogByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	deleteQuery := "DELETE FROM blogs WHERE id = $1"
	rows, err := database.DB.Query(deleteQuery, id)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Blog Deleted Successfully",
	})
}
