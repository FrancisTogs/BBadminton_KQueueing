package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

// 1. Define the data model using a struct and JSON tags
type Player struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level string `json:"level"`
}

// 2. Set up an in-memory database and a mutex to prevent concurrent write crashes
var (
	players = []Player{
		{ID: 1, Name: "Kiko", Level: "Low Advanced"},
		{ID: 2, Name: "Marc", Level: "High Advanced"},
	}
	nextID = 3
	mu     sync.Mutex
)

// 3. Handle GET and POST requests
func playersHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS so your Svelte app on localhost:5173 can access this localhost:8080 API
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle browser preflight checks
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		// Convert Go struct into JSON and send to Svelte
		mu.Lock()
		json.NewEncoder(w).Encode(players)
		mu.Unlock()

	case "POST":
		// Convert incoming JSON from Svelte into a Go struct
		var newPlayer Player
		if err := json.NewDecoder(r.Body).Decode(&newPlayer); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		newPlayer.ID = nextID
		nextID++
		players = append(players, newPlayer)
		mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPlayer)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/api/players", playersHandler)
	println("Go API is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}