package services

import (
	"encoding/json"
	"main/internal/contexter"
	"main/internal/database"
	"main/pkg/response"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

func GetAllBlogComments(w http.ResponseWriter, r *http.Request) {
	blogId := r.PathValue("blogId")

	getCommentsQuery := `
		SELECT id, blog_id, author_id, body FROM comments where blog_id = $1
	`

	rows, err := database.DB.Query(getCommentsQuery, blogId)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	type Comment struct {
		ID       int    `json:"id"`
		Body     string `json:"body"`
		BlogID   int    `json:"blog_id"`
		AuthorID int    `json:"author_id"`
	}

	comments := make([]Comment, 0)
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.BlogID, &c.AuthorID, &c.Body)
		if err != nil {
			response.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		comments = append(comments, c)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string][]Comment{
		"data": comments,
	})
}

func AddBlogComment(w http.ResponseWriter, r *http.Request) {
	blogID := r.PathValue("blogId")

	user, ok := contexter.UserFromContext(r.Context())
	if !ok {
		response.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	type ReqComment struct {
		Body string `json:"body" validate:"required,min=1,max=1000"`
	}
	var comment ReqComment
	json.NewDecoder(r.Body).Decode(&comment)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(&comment)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addBlogQuery := `
		INSERT INTO comments (author_id, blog_id, body, created_at, updated_at)
		VALUES ( (SELECT id FROM users WHERE email = $1 ), $2, $3, $4, $5);
	`
	date := time.Now().UTC()
	rows, err := database.DB.Query(addBlogQuery, user.Email, blogID, comment.Body, date, date)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "comment successfully added to blog",
	})
}

func UpdateBlogComment(w http.ResponseWriter, r *http.Request) {
	blogID := r.PathValue("blogId")
	commentID := r.PathValue("commentId")
	_, ok := contexter.UserFromContext(r.Context())
	if !ok {
		response.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	type ReqComment struct {
		Body string `json:"body" validate:"required,min=1,max=1000"`
	}
	var comment ReqComment
	json.NewDecoder(r.Body).Decode(&comment)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(&comment)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateCommentQuery := "UPDATE comments SET body = $1 WHERE id = $2 and blog_id = $3"
	rows, err := database.DB.Query(updateCommentQuery, comment.Body, commentID, blogID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Comment updated successfully",
	})
}

func DeleteBlogComment(w http.ResponseWriter, r *http.Request) {
	commentID := r.PathValue("commentId")
	if commentID == "" {
		response.Error(w, "comment id param not provided", http.StatusBadRequest)
		return
	}
	_, ok := contexter.UserFromContext(r.Context())
	if !ok {
		response.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	deleteQuery := "DELETE FROM comments WHERE id = $1;"
	rows, err := database.DB.Query(deleteQuery, commentID)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Comment deleted successfully",
	})
}
