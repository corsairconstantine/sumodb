package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Rikishis           RikishiModel
	TournamentsResults TournamentResultModel
	Bouts              BoutModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Rikishis:           RikishiModel{DB: db},
		TournamentsResults: TournamentResultModel{DB: db},
		Bouts:              BoutModel{DB: db},
	}
}
