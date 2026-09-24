package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Player struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Level            string `json:"level"`
	Status           string `json:"status"`
	WaitStartTime    int64  `json:"waitStartTime"`
	TotalGames       int    `json:"totalGames"`
	TotalWins        int    `json:"totalWins"`
	TotalWaitingTime int64  `json:"totalWaitingTime"`
}

type MatchHistory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Score     string `json:"score"`
	StartUnix int64  `json:"startUnix"`
	EndUnix   int64  `json:"endUnix"`
}

type FinishMatchRequest struct {
	Name      string `json:"name"`
	Score     string `json:"score"`
	StartUnix int64  `json:"startUnix"`
	EndUnix   int64  `json:"endUnix"`
	PlayerIDs []int  `json:"playerIds"` // [p1, p2, p3, p4] -> First 2 are Team 1, last 2 are Team 2
}

var db *sql.DB

func bulkPlayersHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}

	if r.Method == http.MethodPost {
		var players []Player
		if err := json.NewDecoder(r.Body).Decode(&players); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		now := time.Now().Unix()
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		stmt, err := tx.Prepare("INSERT INTO players (name, level, status, wait_start_time, total_games, total_wins, total_waiting_time) VALUES (?, ?, 'waiting', ?, 0, 0, 0)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		var insertedPlayers []Player
		for _, p := range players {
			if strings.TrimSpace(p.Name) == "" {
				continue
			}
			result, err := stmt.Exec(strings.TrimSpace(p.Name), p.Level, now)
			if err != nil {
				continue
			}
			id, _ := result.LastInsertId()
			insertedPlayers = append(insertedPlayers, Player{
				ID:               int(id),
				Name:             strings.TrimSpace(p.Name),
				Level:            p.Level,
				Status:           "waiting",
				WaitStartTime:    now,
				TotalGames:       0,
				TotalWins:        0,
				TotalWaitingTime: 0,
			})
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(insertedPlayers)
	}
}

func setCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}
	w.Header().Set("Content-Type", "application/json")
	return false
}

func playersHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, name, level, status, wait_start_time, total_games, total_wins, total_waiting_time FROM players")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		players := []Player{}
		for rows.Next() {
			var p Player
			if err := rows.Scan(&p.ID, &p.Name, &p.Level, &p.Status, &p.WaitStartTime, &p.TotalGames, &p.TotalWins, &p.TotalWaitingTime); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			players = append(players, p)
		}
		json.NewEncoder(w).Encode(players)

	case http.MethodPost:
		var p Player
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now := time.Now().Unix()
		result, err := db.Exec("INSERT INTO players (name, level, status, wait_start_time, total_games, total_wins, total_waiting_time) VALUES (?, ?, ?, ?, 0, 0, 0)",
			p.Name, p.Level, "waiting", now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id, _ := result.LastInsertId()
		p.ID = int(id)
		p.Status = "waiting"
		p.WaitStartTime = now
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)

	case http.MethodPut:
		var p Player
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err := db.Exec("UPDATE players SET name=?, level=?, status=?, wait_start_time=? WHERE id=?",
			p.Name, p.Level, p.Status, p.WaitStartTime, p.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(p)

	case http.MethodDelete:
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		db.Exec("DELETE FROM players WHERE id = ?", id)
		w.WriteHeader(http.StatusOK)
	}
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}

	if r.Method == http.MethodGet {
		rows, err := db.Query("SELECT id, name, score, start_unix, end_unix FROM match_history ORDER BY id ASC")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		history := []MatchHistory{}
		for rows.Next() {
			var h MatchHistory
			if err := rows.Scan(&h.ID, &h.Name, &h.Score, &h.StartUnix, &h.EndUnix); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			history = append(history, h)
		}
		json.NewEncoder(w).Encode(history)
	}
}

func finishMatchHandler(w http.ResponseWriter, r *http.Request) {
	if setCORS(w, r) {
		return
	}

	if r.Method == http.MethodPost {
		var req FinishMatchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 1. Save to Match History Table
		_, err := db.Exec("INSERT INTO match_history (name, score, start_unix, end_unix) VALUES (?, ?, ?, ?)",
			req.Name, req.Score, req.StartUnix, req.EndUnix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Parse scores safely
		t1Score, t2Score := 0, 0
		parts := strings.Split(req.Score, "-")
		if len(parts) == 2 {
			t1Score, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
			t2Score, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
		}

		t1Won := t1Score > t2Score
		t2Won := t2Score > t1Score

		// Update stats for all 4 players
		for i, playerID := range req.PlayerIDs {
			// Index 0 and 1 belong to Team 1, Index 2 and 3 belong to Team 2
			isWinner := false
			if i < 2 && t1Won {
				isWinner = true
			} else if i >= 2 && t2Won {
				isWinner = true
			}

			winIncrement := 0
			if isWinner {
				winIncrement = 1
			}

			_, err = db.Exec(`
				UPDATE players 
				SET total_games = total_games + 1, 
				    total_wins = total_wins + ?,
				    total_waiting_time = total_waiting_time + GREATEST(0, ? - wait_start_time), 
				    wait_start_time = ? 
				WHERE id = ?`,
				winIncrement, req.StartUnix, req.EndUnix, playerID)
			if err != nil {
				log.Println("Error updating player stats:", err)
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/badminton_queue"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize updated players table with total_wins
	db.Exec(`CREATE TABLE IF NOT EXISTS players (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		level VARCHAR(100) NOT NULL,
		status VARCHAR(50) NOT NULL,
		wait_start_time BIGINT NOT NULL,
		total_games INT NOT NULL DEFAULT 0,
		total_wins INT NOT NULL DEFAULT 0,
		total_waiting_time BIGINT NOT NULL DEFAULT 0
	);`)

	db.Exec(`CREATE TABLE IF NOT EXISTS match_history (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		score VARCHAR(50) NOT NULL,
		start_unix BIGINT NOT NULL,
		end_unix BIGINT NOT NULL
	);`)

	http.HandleFunc("/api/players", playersHandler)
	http.HandleFunc("/api/history", historyHandler)
	http.HandleFunc("/api/matches/finish", finishMatchHandler)
	http.HandleFunc("/api/players/bulk", bulkPlayersHandler)

	println("Go MySQL API is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
