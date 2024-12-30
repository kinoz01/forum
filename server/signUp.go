package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

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
	db := OpenDB()
	defer db.Close()
	
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
