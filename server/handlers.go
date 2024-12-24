package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Handle index web page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if CheckHomeRequest(w, r) {
		return
	}
	ParseAndExecute(w, "", "frontend/templates/index.html")
}

// Handle serving both CSS and JS content
func FileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "frontend" + r.URL.Path

	// Read the file from the filesystem
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		ErrorHandler(w, http.StatusForbidden, http.StatusText(http.StatusForbidden), "You don't have permission to access this link!", err)
		return
	}

	// Serve the file content
	http.ServeContent(w, r, filePath, time.Now(), bytes.NewReader(fileBytes))
}

// signUpHandler expects JSON { "email": "...", "username": "...", "password": "..." }
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Basic validation
	if user.Email == "" || user.Username == "" || user.Password == "" {
		http.Error(w, "Email, username, and password are required", http.StatusBadRequest)
		return
	}
	fmt.Println(user.Email, user.Username, user.Password)


	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert user into DB
	insertUser := `INSERT INTO users (email, username, password) VALUES (?, ?, ?)`
	_, err = db.Exec(insertUser, user.Email, user.Username, hashedPassword)
	if err != nil {
		// Could be unique constraint error, etc.
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Success
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
