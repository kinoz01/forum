package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/mattn/go-sqlite3" 
)

// Initialise ports (api & application ports) using environement variables.
// For example use export PORT=:$PORT to set port where the server should start.
func (p *Ports) InitialisePorts() {
	p.Port = os.Getenv("PORT")
	p.ApiPort = os.Getenv("APIPORT")
}

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

// Initialise the APILinks global struct using apiLinks.json.
func InitialiseApiLinks() {
	// Load the JSON with needed links.
	var err error
	APILinks, err = LoadApiLinks("./server/apiLinks.json")
	if err != nil {
		log.Fatalf("Error loading API links: %v", err)
	}
}

func LoadApiLinks(filename string) (*ApiLinks, error) {
	// Load the JSON file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into ApiLinks struct
	var links ApiLinks
	err = json.Unmarshal(content, &links)
	if err != nil {
		return nil, err
	}

	return &links, nil
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

func InitDB() *sql.DB {
	var err error
	// Open or create SQLite DB
	db, err = sql.Open("sqlite3", "./dataBase/forum.db")
	if err != nil {
		log.Fatal("Failed to open SQLite database:", err)
	}

	// Create users table if not exists
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT NOT NULL UNIQUE,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	if _, err := db.Exec(createUsersTable); err != nil {
		log.Fatal("Failed to create users table:", err)
	}
	return db
}
