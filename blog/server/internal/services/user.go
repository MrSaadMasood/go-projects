package services

import (
	"encoding/json"
	"main/internal/contexter"
	"main/internal/database"
	"main/pkg/response"
	"net/http"
)

func UserProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := contexter.UserFromContext(r.Context())
	if !ok {
		response.Error(w, "user must be logged in to get its profile", http.StatusUnauthorized)
		return
	}

	userBlogListQuery := "SELECT title FROM blogs WHERE author_id = ( SELECT id FROM users WHERE EMAIL = $1)"
	rows, err := database.DB.Query(userBlogListQuery, user.Email)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()
	blogs := make([]string, 0)

	for rows.Next() {
		var title string
		err := rows.Scan(&title)
		if err != nil {
			response.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		blogs = append(blogs, title)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data": map[string]any{
			"userProfile": map[string]any{
				"email": user.Email,
				"blogs": blogs,
			},
		},
	})
}
