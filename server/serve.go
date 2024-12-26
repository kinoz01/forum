package server

import (
	"log"
	"net"
	"net/http"
)

func Serve() {
	Initialise()
	Route()
	Listen()
}

// Define server endPoints.
func Route() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/css/", FileHandler)
	http.HandleFunc("/js/", FileHandler)
	http.HandleFunc("/signup", SignUpHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/session-check", SessionCheckHandler)
	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetPostsHandler(w, r)
		case http.MethodPost:
			CreatePostHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// Listen on a env Specific or random PORT for incoming requests.
func Listen() {
	listener, err := net.Listen("tcp", ":"+Port.Port)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
		return
	}
	_, p, _ := net.SplitHostPort(listener.Addr().String())

	log.Printf("Starting server at http://127.0.0.1:%s", p)
	if err = http.Serve(listener, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
