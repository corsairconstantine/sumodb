package apiserver

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
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
	if r.URL.Query() != nil {
		h.rikishisQuery(w, r)
	} else {
		switch r.Method {
		case http.MethodGet:
			h.getRikishis(w, r)
		case http.MethodPost:
			h.addRikishi(w, r)
		default:
			http.Error(w, "Wrong method", 405)
		}
	}
}

func (h *rikishisHandler) getRikishis(w http.ResponseWriter, r *http.Request) {
	db := h.db
	rows, err := db.Query("SELECT * FROM rikishi;")
	if err != nil {
		h.logger.Println(err)
	}
	var rikishis []rikishi
	for rows.Next() {
		var r rikishi
		if err := rows.Scan(&r.ID, &r.Shikona, &r.Rank, &r.Height, &r.Weight); err != nil {
			h.logger.Println(err)
		}
		rikishis = append(rikishis, r)
	}
	json.NewEncoder(w).Encode(rikishis)
}

func (h *rikishisHandler) addRikishi(w http.ResponseWriter, r *http.Request) {
	db := h.db
	var newrikishi rikishi
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newrikishi)
	if err != nil {
		h.logger.Println(err)
	}
	result, err := db.Exec("INSERT INTO rikishi (shikona, rank, height, weight) VALUES ($1, $2, $3, $4);",
		newrikishi.Shikona, newrikishi.Rank, newrikishi.Height, newrikishi.Weight)
	if err != nil {
		h.logger.Println(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		h.logger.Println(err)
	}
	h.logger.Printf("Added a row with id: %v\n", id)
}

func (h *rikishisHandler) rikishisQuery(w http.ResponseWriter, r *http.Request) {

}
