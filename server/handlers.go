package server

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Handle index web page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if CheckHomeRequest(w, r) {
		return
	}
	ParseAndExecute(w, "", "static/templates/index.html")
}

// Handle serving static content.
func FileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "static" + r.URL.Path

	// Read the file from the filesystem
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		ErrorHandler(w, http.StatusForbidden, http.StatusText(http.StatusForbidden), "You don't have permission to access this link!", err)
		return
	}
	// Serve the file content
	http.ServeContent(w, r, filePath, time.Now(), bytes.NewReader(fileBytes))
}

// sessionCheckHandler checks whether the user has a valid session.
func SessionCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Try to get the user from the session.
	user, err := getUserFromSession(r) // a function that verifies session_token cookie in DB
	if err != nil {
		// Not logged in or session invalid
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"loggedIn": false}`)
		return
	}

	// If user is logged in, return username (+ profilePic if you have it).
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"loggedIn": true, "username": %q}`, user.Username)
}

func getUserFromSession(r *http.Request) (*User, error) {
	// Get the session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, fmt.Errorf("no session token provided")
	}

	token := cookie.Value
	var session struct {
		UserID    int
		ExpiresAt time.Time
	}

	// Open the DB and query the sessions table
	db := openDB()
	defer db.Close()

	err = db.QueryRow(`SELECT user_id, expires_at FROM sessions WHERE token = ?`, token).
		Scan(&session.UserID, &session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired session token")
	}

	// Check if the session is expired
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	// Fetch the user associated with the session
	var user User
	err = db.QueryRow(`SELECT id, email, username FROM users WHERE id = ?`, session.UserID).
		Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}
