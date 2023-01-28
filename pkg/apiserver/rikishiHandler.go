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
		if err := h.getRikishi(w, r); err != nil {
			json.NewEncoder(w).Encode(err)
		}
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (h *rikishiHandler) getRikishi(w http.ResponseWriter, r *http.Request) *appError {
	param := strings.TrimPrefix(r.URL.Path, "/rikishis/id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return &appError{err, "Bad request, could not convert parameter to id", 400}
	}

	row := h.db.QueryRow("SELECT * FROM rikishi WHERE id = $1", id)
	var rik rikishi
	if err = row.Scan(&rik.ID, &rik.Shikona, &rik.Rank, &rik.Height, &rik.Weight); err != nil {
		return &appError{err, "no data found", 404}
	}
	json.NewEncoder(w).Encode(&rik)
	return nil
}
