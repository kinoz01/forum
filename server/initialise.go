package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
)

func Initialise() {
	Port.InitialisePorts()
	InitialiseLinks()
	db := InitialiseDB()
	defer db.Close()
}

// Initialise ports (api & application ports) using environement variables.
// For example use export PORT=:$PORT to set port where the server should start.
func (p *Ports) InitialisePorts() {
	p.Port = os.Getenv("PORT")
	p.ApiPort = os.Getenv("APIPORT")
}

// Initialise the Links global struct using apiLinks.json.
func InitialiseLinks() {
	// Load the JSON with needed links.
	var err error
	APILinks, err = LoadLinks("./server/apiLinks.json")
	if err != nil {
		log.Fatalf("Error loading API links: %v", err)
	}
}

// Open the Json file and unmarshal it's content into Links struct. 
func LoadLinks(filename string) (*Links, error) {
	file, content, err := OpenRead(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Unmarshal the JSON into ApiLinks struct
	var links Links
	err = json.Unmarshal(content, &links)
	if err != nil {
		return nil, err
	}
	return &links, nil
}

// Create/open DB and create tables.
func InitialiseDB() *sql.DB {
	// Open or create SQLite DB
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal("Failed to open SQLite database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("connection to the database is dead.", err)
	}

	file, content, err := OpenRead("./database/schema.sql")
	if err != nil {
		log.Fatal("Failed to get database tables:", err)
	}
	defer file.Close()
    
	if _, err := db.Exec(string(content)); err != nil {
		log.Fatal("Failed to create database tables:", err)
	}
	return db
}
