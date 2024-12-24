package server

import "database/sql"

type Ports struct {
	Port    string
	ApiPort string
}

// Declare global variables Ports and ApiLinks.
var (
	Port     Ports
	APILinks *ApiLinks
)

// ApiLinks represents the JSON links structure
type ApiLinks struct {
	ErrorPage string `json:"error"`
}

// User is a simple struct to capture signup form data.
type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// DB global variable for simplicity (avoid in production).
var db *sql.DB
