package server

import (
	"encoding/json"
	"net/http"
	"time"
)

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db := OpenDB()
	defer db.Close()

	rows, err := db.Query(`
        SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username
        FROM posts p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC
    `)
	if err != nil {
		http.Error(w, "Failed to query posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// We'll store a "PostWithAuthor" to also return the username
	type PostWithAuthor struct {
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		Username  string    `json:"username"`
		Title     string    `json:"title"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}

	var posts []PostWithAuthor
	for rows.Next() {
		var pa PostWithAuthor
		if err := rows.Scan(&pa.ID, &pa.UserID, &pa.Title, &pa.Content, &pa.CreatedAt, &pa.Username); err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, pa)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Verify user is logged in (session)
    user, err := getUserFromSession(r)
    if err != nil {
        http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
        return
    }

    // Parse JSON
    var payload struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        return
    }

    // Insert post
    db := OpenDB()
    defer db.Close()

    _, err = db.Exec(`
        INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`,
        user.ID, payload.Title, payload.Content)
    if err != nil {
        http.Error(w, "Failed to create post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Post created successfully"))
}
