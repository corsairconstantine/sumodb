package apiserver

import (
	"encoding/json"
	"net/http"
)

func (s *APIserver) rikishiHandler(w http.ResponseWriter, r *http.Request) {
	db := s.store.Db
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	//get all rikishis
	case http.MethodGet:
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
	//add a rikishi
	case http.MethodPost:
		var rik rikishi
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&rik)
		if err != nil {
			s.logger.Panic(err)
		}
		s.logger.Printf("Inserting data: shikona:%s, rank: %s, height: %d, weight: %d\n", rik.Shikona, rik.Rank, rik.Height, rik.Weight)
		_, err = db.Query("INSERT INTO rikishi (shikona, rank, height, weight) VALUES ($1, $2, $3, $4);", rik.Shikona, rik.Rank, rik.Height, rik.Weight)
		if err != nil {
			s.logger.Panic(err)
		}
	//update a rikishi
	case http.MethodPut:
		var rik rikishi
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&rik)
		if err != nil {
			s.logger.Panic(err)
		}
		s.logger.Printf("Updating data: shikona:%s, rank: %s, height: %d, weight: %d\n", rik.Shikona, rik.Rank, rik.Height, rik.Weight)
		_, err = db.Query("UPDATE rikishi SET shikona = $2, rank = $3, height = $4, weight = $5 WHERE id = $1;", rik.ID, rik.Shikona, rik.Rank, rik.Height, rik.Weight)
		if err != nil {
			s.logger.Panic(err)
		}
	}
}
