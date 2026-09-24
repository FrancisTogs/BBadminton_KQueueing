# Badminton Queue & Match Manager

A full-stack web application designed for recreational badminton groups, clubs, or facilities to manage player registrations, skill-tiered waiting queues, active court rotations, live match timers, and historical performance tracking.

---

## 🛠️ Tech Stack

* **Frontend:** Svelte 5 (TypeScript), HTML5, CSS3, Browser `localStorage` caching for live session states.
* **Backend:** Go (Golang) using `net/http` and standard database drivers.
* **Database:** MySQL (handles permanent player statistics and match history records).

---

## ✨ Features

1. **Active Court Management:** Dynamically add or remove courts, monitor live match durations with real-time tickers, and conclude matches with final scores[cite: 15].
2. **Skill-Tiered Player Management:** Register players across various skill levels (Beginner, Intermediate, Advanced)[cite: 15]. Track their historical total games, total wins, and accumulated waiting times.
3. **Smart Queue System:** Select 4 available players to form Team 1 and Team 2, queue matches, edit lineups on the fly, and push matches directly to open courts.
4. **Match History & Winner Highlighting:** Automatically calculates the winning team based on entered final scores, displays crown indicators (`👑`), highlights winning team boxes in green, and records overall match timelines.
5. **Robust Two-Tier State Architecture:** Live session data (active courts and queues) is instantly cached in browser storage to survive page refreshes, while critical records (player profiles and match outcomes) persist safely in MySQL via the Go backend.

---

## 🚀 Getting Started

### Prerequisites
* **Go** (v1.18 or higher recommended)
* **Node.js** & npm (or a compatible package manager like Vite/Bun)
* **MySQL Server** running locally (`127.0.0.1:3306`)

---

### Database Setup

1. Create a MySQL database named `badminton_queue`:

CREATE DATABASE badminton_queue;

'(Note: The Go backend will automatically initialize the required players and match_history tables upon startup.)'

2. Run the Go Backend
Navigate to your backend directory where main.go is located, and start the server:

go run main.go

The API server will launch and listen on http://localhost:8080.

3. Run the Svelte Frontend
Navigate to your frontend project directory, install dependencies, and start the development server:

# Install dependencies
npm install# Start development server
npm run dev
Open your browser and navigate to the local development URL provided by Vite (typically http://localhost:5173) to start managing your badminton queue!
