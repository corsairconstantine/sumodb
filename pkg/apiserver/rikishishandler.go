package apiserver

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type rikishi struct {
	ID      int    `json:"id"`
	Shikona string `json:"shikona"`
	Rank    string `json:"rank"`
	Height  uint8  `json:"height"`
	Weight  uint8  `json:"weight"`
}

type rikishisHandler struct {
	logger *log.Logger
	db     *sql.DB
}

func (h *rikishisHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.URL.RawQuery != "" {
		if err := h.rikishisQuery(w, r); err != nil {
			json.NewEncoder(w).Encode(err)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			if err := h.getRikishis(w, r); err != nil {
				json.NewEncoder(w).Encode(err)
			}
		case http.MethodPost:
			if err := h.addRikishi(w, r); err != nil {
				json.NewEncoder(w).Encode(err)
			}
		default:
			http.Error(w, "Wrong method", 405)
		}
	}
}

func (h *rikishisHandler) getRikishis(w http.ResponseWriter, r *http.Request) *appError {
	rows, err := h.db.Query("SELECT * FROM rikishi;")
	if err != nil {
		return &appError{err, err.Error(), 500}
	}
	var rikishis []rikishi
	for rows.Next() {
		var r rikishi
		if err := rows.Scan(&r.ID, &r.Shikona, &r.Rank, &r.Height, &r.Weight); err != nil {
			return &appError{err, err.Error(), 404}
		}
		rikishis = append(rikishis, r)
	}
	json.NewEncoder(w).Encode(&rikishis)
	return nil
}

func (h *rikishisHandler) addRikishi(w http.ResponseWriter, r *http.Request) *appError {
	db := h.db
	var newrikishi rikishi
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newrikishi)
	if err != nil {
		return &appError{err, err.Error(), 400}
	}
	result, err := db.Exec("INSERT INTO rikishi (shikona, rank, height, weight) VALUES ($1, $2, $3, $4);",
		newrikishi.Shikona, newrikishi.Rank, newrikishi.Height, newrikishi.Weight)
	if err != nil {
		return &appError{err, err.Error(), 404}
	}
	id, err := result.LastInsertId()
	if err != nil {
		return &appError{err, err.Error(), 500}
	}
	h.logger.Printf("Added a row with id: %v\n", id) //return this message to the client
	return nil
}

func (h *rikishisHandler) rikishisQuery(w http.ResponseWriter, r *http.Request) *appError {
	db := h.db
	vals, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return &appError{err, "Invalid query", 400}
	}
	rows, err := db.Query("SELECT * FROM rikishi WHERE shikona ILIKE '%' || $1 || '%' AND rank ILIKE '%' || $2 || '%';",
		vals.Get("shikona"), vals.Get("rank"))
	if err != nil {
		return &appError{err, err.Error(), 500}
	}
	defer rows.Close()
	var riks []rikishi
	for rows.Next() {
		var rik rikishi
		err := rows.Scan(&rik.ID, &rik.Shikona, &rik.Rank, &rik.Height, &rik.Weight)
		if err != nil {
			return &appError{err, err.Error(), 500}
		}
		riks = append(riks, rik)
	}
	json.NewEncoder(w).Encode(&riks)
	return nil
}
