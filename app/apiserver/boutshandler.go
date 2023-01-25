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
		rows, err := db.Query("SELECT W.id, W.id, W.shikona, W.rank, W.height, W.weight, loser, rikishi.shikona, rikishi.rank, rikishi.height, rikishi.weight, tournament, division, day FROM (SELECT bout.id, winner, shikona, rank, height, weight, loser, tournament, division, day FROM bout INNER JOIN rikishi ON winner = rikishi.id) AS W INNER JOIN rikishi ON loser = rikishi.id;")
		if err != nil {
			s.logger.Fatal(err)
		}
		defer rows.Close()
		var bouts []bout
		for rows.Next() {
			var b bout
			err := rows.Scan(&b.ID, &b.Winner.ID, &b.Winner.Shikona, &b.Winner.Rank, &b.Winner.Height, &b.Winner.Weight,
				&b.Loser.ID, &b.Loser.Shikona, &b.Loser.Rank, &b.Loser.Height, &b.Loser.Weight,
				&b.Tournament, &b.Division, &b.Day)
			if err != nil {
				s.logger.Fatal(err)
			}
			bouts = append(bouts, b)
		}
		json.NewEncoder(w).Encode(&bouts)
		// adding a bout
	case http.MethodPost:
		var newbout bout
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newbout)
		if err != nil {
			s.logger.Panic(err)
		} // ---TO DO probably should check if rikishi ids exist in the db UPDATE nvm postgres does it for me
		s.logger.Printf("Inserting data: winner:%v, loser: %v, tournament: %s, division: %s, day: %v\n", newbout.Winner, newbout.Loser, newbout.Tournament, newbout.Division, newbout.Day)
		_, err = db.Query("INSERT INTO bout (winner, loser, tournament, division, day) VALUES ($1, $2, $3, $4, $5);", newbout.Winner, newbout.Loser, newbout.Tournament, newbout.Division, newbout.Day)
		if err != nil {
			s.logger.Panic(err)
		}
	// updating a bout
	case http.MethodPut:
		var newbout bout
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newbout)
		if err != nil {
			s.logger.Panic(err)
		} // ---TO DO probably should check if rikishi ids exist in the db UPDATE nvm postgres does it for me
		_, err = db.Query("UPDATE bout SET winner = $2, loser = $3, tournament = $4, division = $5, day = $6 WHERE id = $1;", newbout.ID, newbout.Winner, newbout.Loser, newbout.Tournament, newbout.Division, newbout.Day)
		if err != nil {
			s.logger.Panic(err)
		}
		s.logger.Printf("Updated bout: bout id: %v, winner:%v, loser: %v, tournament: %s, division: %s, day: %v\n", newbout.ID, newbout.Winner, newbout.Loser, newbout.Tournament, newbout.Division, newbout.Day)
		/*			//deleting a bout
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
