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
	case http.MethodPut:
		if err := h.updateRikishi(w, r); err != nil {
			json.NewEncoder(w).Encode(err)
		}
	case http.MethodDelete:
		if err := h.deleteRikishi(w, r); err != nil {
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

func (h *rikishiHandler) updateRikishi(w http.ResponseWriter, r *http.Request) *appError {
	var rik rikishi
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rik)
	if err != nil {
		return &appError{err, err.Error(), 400}
	}
	_, err = h.db.Exec("UPDATE rikishi SET shikona = $2, rank = $3, height = $4, weight = $5 WHERE id = $1",
		rik.ID, rik.Shikona, rik.Rank, rik.Height, rik.Weight)
	if err != nil {
		return &appError{err, err.Error(), 404}
	}
	return nil
}

func (h *rikishiHandler) deleteRikishi(w http.ResponseWriter, r *http.Request) *appError {
	param := strings.TrimPrefix(r.URL.Path, "/rikishis/id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return &appError{err, "Bad request, could not convert parameter to id", 400}
	}

	_, err = h.db.Exec("DELETE FROM rikishi WHERE id = $1", id)
	if err != nil {
		return &appError{err, err.Error(), 404}
	}
	json.NewEncoder(w).Encode("Done")
	return nil
}
