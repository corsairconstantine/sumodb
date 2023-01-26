package apiserver

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type bout struct {
	ID         int     `json:"id"`
	Winner     rikishi `json:"winner"`
	Loser      rikishi `json:"loser"`
	Tournament string  `json:"tournament"`
	Division   string  `json:"division"`
	Day        uint8   `json:"day"`
}

type boutsHandler struct {
	logger *log.Logger
	db     *sql.DB
}

func (h *boutsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Query() != nil {
		h.boutsQuery(w, r)
	} else {
		switch r.Method {
		case http.MethodGet:
			h.getBouts(w, r)
		case http.MethodPost:
			h.addBout(w, r)
		default:
			http.Error(w, "Wrong method", 405)
		}
	}
}

func (h *boutsHandler) getBouts(w http.ResponseWriter, r *http.Request) {
	db := h.db
	rows, err := db.Query("SELECT W.id, W.id, W.shikona, W.rank, W.height, W.weight, loser, rikishi.shikona, rikishi.rank, rikishi.height, rikishi.weight, tournament, division, day FROM (SELECT bout.id, winner, shikona, rank, height, weight, loser, tournament, division, day FROM bout INNER JOIN rikishi ON winner = rikishi.id) AS W INNER JOIN rikishi ON loser = rikishi.id;")
	if err != nil {
		h.logger.Fatal(err)
	}
	defer rows.Close()
	var bouts []bout
	for rows.Next() {
		var b bout
		err := rows.Scan(&b.ID, &b.Winner.ID, &b.Winner.Shikona, &b.Winner.Rank, &b.Winner.Height, &b.Winner.Weight,
			&b.Loser.ID, &b.Loser.Shikona, &b.Loser.Rank, &b.Loser.Height, &b.Loser.Weight,
			&b.Tournament, &b.Division, &b.Day)
		if err != nil {
			h.logger.Fatal(err)
		}
		bouts = append(bouts, b)
	}
	json.NewEncoder(w).Encode(&bouts)
}

func (h *boutsHandler) addBout(w http.ResponseWriter, r *http.Request) {
	db := h.db
	var newbout bout
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newbout)
	if err != nil {
		h.logger.Println(err)
	}
	result, err := db.Exec("INSERT INTO bout (winner, loser, tournament, division, day) VALUES ($1, $2, $3, $4, $5);",
		newbout.Winner, newbout.Loser, newbout.Tournament, newbout.Division, newbout.Day)
	if err != nil {
		h.logger.Println(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		h.logger.Println(err)
	}
	h.logger.Printf("Added a bout at id: %v\n", id)
}

func (h *boutsHandler) boutsQuery(w http.ResponseWriter, r *http.Request) {

}
