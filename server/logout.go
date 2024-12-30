package server

import (
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "No session token provided", http.StatusUnauthorized)
		return
	}

	token := cookie.Value

	db := OpenDB()
	defer db.Close()

	// Remove session from DB
	_, err = db.Exec(`DELETE FROM sessions WHERE token = ?`, token)
	if err != nil {
		http.Error(w, "Failed to remove session", http.StatusInternalServerError)
		return
	}

	// Expire the cookie
	expiredCookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, expiredCookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
