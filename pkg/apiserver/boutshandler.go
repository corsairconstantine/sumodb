package apiserver

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type bout struct {
	ID     int     `json:"id"`
	Winner rikishi `json:"winner"`
	Loser  rikishi `json:"loser"`
	//WinnerRank string `json:"winnerrank"`
	//LoserRank string `json:"loserrank"`
	Tournament string `json:"tournament"`
	Division   string `json:"division"`
	Day        uint8  `json:"day"`
}

type boutsHandler struct {
	logger *log.Logger
	db     *sql.DB
}

func (h *boutsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.URL.RawQuery != "" {
		if err := h.boutsQuery(w, r); err != nil {
			json.NewEncoder(w).Encode(err)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			if err := h.getBouts(w, r); err != nil {
				json.NewEncoder(w).Encode(&err)
			}
		case http.MethodPost:
			if err := h.addBout(w, r); err != nil {
				json.NewEncoder(w).Encode(&err)
			}
		default:
			http.Error(w, "Wrong method", 405)
		}
	}
}

func (h *boutsHandler) getBouts(w http.ResponseWriter, r *http.Request) *appError {
	rows, err := h.db.Query("SELECT W.id, winnerId, W.shikona, W.rank, W.height, W.weight, loser, rikishi.shikona, rikishi.rank, rikishi.height, rikishi.weight, tournament, division, day FROM (SELECT bout.id, winner AS winnerID, shikona, rank, height, weight, loser, tournament, division, day FROM bout INNER JOIN rikishi ON winner = rikishi.id) AS W INNER JOIN rikishi ON loser = rikishi.id;")
	if err != nil {
		return &appError{err, err.Error(), 404}
	}
	defer rows.Close()
	var bouts []bout
	for rows.Next() {
		var b bout
		err := rows.Scan(&b.ID, &b.Winner.ID, &b.Winner.Shikona, &b.Winner.Rank, &b.Winner.Height, &b.Winner.Weight,
			&b.Loser.ID, &b.Loser.Shikona, &b.Loser.Rank, &b.Loser.Height, &b.Loser.Weight,
			&b.Tournament, &b.Division, &b.Day)
		if err != nil {
			return &appError{err, err.Error(), 500}
		}
		bouts = append(bouts, b)
	}
	json.NewEncoder(w).Encode(&bouts)
	return nil
}

func (h *boutsHandler) addBout(w http.ResponseWriter, r *http.Request) *appError {
	var newbout bout
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newbout)
	if err != nil {
		return &appError{err, err.Error(), 400}
	}
	result, err := h.db.Exec("INSERT INTO bout (winner, loser, tournament, division, day) VALUES ($1, $2, $3, $4, $5);",
		newbout.Winner, newbout.Loser, newbout.Tournament, newbout.Division, newbout.Day)
	if err != nil {
		return &appError{err, err.Error(), 404}
	}
	id, err := result.LastInsertId()
	if err != nil {
		return &appError{err, err.Error(), 404}
	}
	h.logger.Printf("Added a bout at id: %v\n", id)
	return nil
}

func (h *boutsHandler) boutsQuery(w http.ResponseWriter, r *http.Request) *appError {
	vals, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return &appError{err, "Invalid query", 400}
	}
	rows, err := h.db.Query("SELECT W.id, winnerId, W.shikona, W.rank, W.height, W.weight, loser, rikishi.shikona, rikishi.rank, rikishi.height, rikishi.weight, tournament, division, day FROM (SELECT bout.id , winner AS winnerId, shikona, rank, height, weight, loser, tournament, division, day FROM bout INNER JOIN rikishi ON winner = rikishi.id) AS W INNER JOIN rikishi ON loser = rikishi.id WHERE W.shikona ILIKE '%' || $1 || '%' AND rikishi.shikona ILIKE '%' || $2 || '%' AND tournament ILIKE '%' || $3 || '%';",
		vals.Get("winner"), vals.Get("loser"), vals.Get("tournament"))
	if err != nil {
		return &appError{err, err.Error(), 404}
	}
	defer rows.Close()
	var bouts []bout
	for rows.Next() {
		var b bout
		err := rows.Scan(&b.ID, &b.Winner.ID, &b.Winner.Shikona, &b.Winner.Rank, &b.Winner.Height, &b.Winner.Weight,
			&b.Loser.ID, &b.Loser.Shikona, &b.Loser.Rank, &b.Loser.Height, &b.Loser.Weight,
			&b.Tournament, &b.Division, &b.Day)
		if err != nil {
			return &appError{err, err.Error(), 500}
		}
		bouts = append(bouts, b)
	}
	json.NewEncoder(w).Encode(&bouts)
	return nil
}
