package apiserver

import (
	"encoding/json"
	"net/http"
)

type rikishi struct {
	ID      int    `json:"id"`
	Shikona string `json:"shikona"`
	Rank    string `json:"rank"`
	Height  uint8  `json:"height"`
	Weight  uint8  `json:"weight"`
}

func (s *APIserver) rikishiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		s.getRikishis(w, r)
	case http.MethodPost:
		s.addRikishi(w, r)
	default:
		http.Error(w, "Wrong method", 405)
	}
}

func (s *APIserver) getRikishis(w http.ResponseWriter, r *http.Request) {
	db := s.store.Db
	rows, err := db.Query("SELECT * FROM rikishi;")
	if err != nil {
		s.logger.Fatal(err)
	}
	var rikishis []rikishi
	for rows.Next() {
		var r rikishi
		if err := rows.Scan(&r.ID, &r.Shikona, &r.Rank, &r.Height, &r.Weight); err != nil {
			s.logger.Fatal(err)
		}
		rikishis = append(rikishis, r)
	}
	json.NewEncoder(w).Encode(rikishis)
}

func (s *APIserver) addRikishi(w http.ResponseWriter, r *http.Request) {
	db := s.store.Db
	var newrikishi rikishi
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newrikishi)
	if err != nil {
		s.logger.Println(err)
	}
	result, err := db.Exec("INSERT INTO rikishi (shikona, rank, height, weight) VALUES ($1, $2, $3, $4);",
		newrikishi.Shikona, newrikishi.Rank, newrikishi.Height, newrikishi.Weight)
	if err != nil {
		s.logger.Println(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		s.logger.Println(err)
	}
	s.logger.Printf("Added a row with id: %v\n", id)
}
