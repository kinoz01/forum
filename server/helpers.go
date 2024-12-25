package server

import (
	"bytes"
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func CheckHomeRequest(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path != "/" {
		ErrorHandler(w, 404, "Look like you're lost!", "The page you are looking for is not available!", nil)
		return true
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, 405, http.StatusText(http.StatusMethodNotAllowed), "Only GET method is allowed!", nil)
		return true
	}
	return false
}

// Parse the html files and excute them after checking for errors.
func ParseAndExecute(w http.ResponseWriter, data any, file string) {
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, "Something seems wrong, try again later!", "Internal Server Error!", err)
		return
	}

	// write to a temporary buffer instead of writing directly to w.
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		ErrorHandler(w, 500, "Something seems wrong, try again later!", "Internal Server Error!", err)
		return
	}
	// If successful, write the buffer content to the ResponseWriter
	buf.WriteTo(w)
}

func openDB() *sql.DB {
	// Open or create SQLite DB
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal("Failed to open SQLite database:", err)
	}
	return db
}

func OpenRead(filename string) (*os.File, []byte, error) {
	// Load the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}
	return file, content, nil
}
