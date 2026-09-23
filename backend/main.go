package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Player struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level string `json:"level"`
}

var db *sql.DB

func playersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Added PUT to allowed methods
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, name, level FROM players")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		players := []Player{}
		for rows.Next() {
			var p Player
			if err := rows.Scan(&p.ID, &p.Name, &p.Level); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			players = append(players, p)
		}
		json.NewEncoder(w).Encode(players)

	case http.MethodPost:
		var newPlayer Player
		if err := json.NewDecoder(r.Body).Decode(&newPlayer); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := db.Exec("INSERT INTO players (name, level) VALUES (?, ?)", newPlayer.Name, newPlayer.Level)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()
		newPlayer.ID = int(id)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPlayer)

	case http.MethodPut:
		// Decode the updated player data
		var updatedPlayer Player
		if err := json.NewDecoder(r.Body).Decode(&updatedPlayer); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update the database record matching the ID
		_, err := db.Exec("UPDATE players SET name = ?, level = ? WHERE id = ?", updatedPlayer.Name, updatedPlayer.Level, updatedPlayer.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedPlayer)

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM players WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	var err error
	
	// Make sure this DSN matches the one that successfully connected earlier
	dsn := "root:@tcp(127.0.0.1:3306)/badminton_queue"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS players (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		level VARCHAR(100) NOT NULL
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/players", playersHandler)
	println("Go MySQL API is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}