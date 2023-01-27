package apiserver

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type rikishiHandler struct {
	logger *log.Logger
	db     *sql.DB
}

func (h *rikishiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		h.getRikishi(w, r)
	default:
		http.Error(w, "Wrong method", 405)
	}
}

func (h *rikishiHandler) getRikishi(w http.ResponseWriter, r *http.Request) {
	db := h.db
	param := strings.TrimPrefix(r.URL.Path, "/rikishis/id")
	id, err := strconv.Atoi(param)
	if err != nil {
		h.logger.Println(err)
		http.Error(w, "Bad request. Could not convert parameter to id", 400)
	} else {
		row := db.QueryRow("SELECT * FROM rikishi WHERE id = $1", id)
		var rik rikishi
		if err = row.Scan(&rik.ID, &rik.Shikona, &rik.Rank, &rik.Height, &rik.Weight); err != nil {
			h.logger.Println(err)
		}
		json.NewEncoder(w).Encode(rik)
	}
}
