package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find user by email
	var user User
	db := OpenDB()
	defer db.Close()

	row := db.QueryRow(`SELECT id, email, username, password FROM users WHERE email = ?`, creds.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		row := db.QueryRow(`SELECT id, email, username, password FROM users WHERE username = ?`, creds.Email)
		err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
		if err != nil {
			http.Error(w, "User not found or invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	// Compare password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session token
	token := uuid.New().String()

	// Set expiration, e.g., 24 hours from now
	expiresAt := time.Now().Add(24 * time.Hour)

	// Insert session into DB
	_, err = db.Exec(`INSERT INTO sessions (user_id, token, expires_at) VALUES (?, ?, ?)`,
		user.ID, token, expiresAt)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set the token in a cookie
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  expiresAt,
		HttpOnly: true, // Helps mitigate XSS
		Path:     "/",  // Cookie valid for entire site
	}
	http.SetCookie(w, cookie)

	// Return JSON with user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Login successful",
		"username": user.Username,
	})
}
