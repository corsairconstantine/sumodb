package apiserver

import (
	"encoding/json"
	"net/http"
)

/*
bouts handler
cases:
get all bouts
get bouts by rikishi's name
get bouts of a specific tournament
get bouts rikishi vs rikishi
*/
func (s *APIserver) boutsHandler(w http.ResponseWriter, r *http.Request) {
	db := s.store.Db
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		//getting all bouts
		rows, err := db.Query("SELECT q.id, shikona, loser, tournament, division, day FROM (SELECT bout.id, winner, shikona as loser, tournament, division, day FROM bout INNER JOIN rikishi ON loser = rikishi.id) AS q INNER JOIN rikishi ON winner = rikishi.id;")
		if err != nil {
			s.logger.Fatal(err)
		}
		var bouts []bout
		for rows.Next() {
			var b bout
			if err := rows.Scan(&b.ID, &b.Winner, &b.Loser, &b.Tournament, &b.Division, &b.Day); err != nil {
				s.logger.Fatal(err)
			}
			bouts = append(bouts, b)
		}
		json.NewEncoder(w).Encode(&bouts)
		/*	// adding a bout
			case http.MethodPost:
				result, err := db.Query("INSERT INTO bouts (winner, loser, tournament, division, day) VALUES")
					// updating a bout
					case http.MethodPut:
						for _, item := range bouts {
							if item.ID == id {
								json.NewEncoder(w).Encode(item)
							}
						}
					//deleting a bout
					case http.MethodDelete:
						for index, item := range bouts {
							if item.ID == id {
								bouts = append(bouts[:index], bouts[index+1:]...)
								break
							}
						}
						json.NewEncoder(w).Encode(bouts)*/
	}
}
