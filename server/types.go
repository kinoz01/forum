package server

import "time"

// Global variables Ports and ApiLinks.
var (
	Port     Ports
	APILinks *Links
)

type Ports struct {
	Port    string
	ApiPort string
}

// Links represents the JSON links structure
type Links struct {
	ErrorPage string `json:"error"`
}

// User is a simple struct to capture signup form data.
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
