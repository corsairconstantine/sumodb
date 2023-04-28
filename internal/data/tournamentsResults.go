package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/corsairconstantine/sumodb/internal/validator"
	"github.com/lib/pq"
)

type TournamentResult struct {
	ID         int64  `json:"id"`
	Tournament string `json:"tournament"`
	Rikishi    string `json:"rikishi"`
	Rank       string `json:"rank"`
	Wins       int32  `json:"wins"`
	Losses     int32  `json:"losses"`
	Absent     int32  `json:"Absent"`
	Version    int32  `json:"version"`
}

type TournamentResultModel struct {
	DB *sql.DB
}

func (t TournamentResultModel) Insert(tr *TournamentResult) error {
	query := `
		INSERT INTO tournaments_results (tournament, rikishi, rank, wins, losses, absent)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, version`

	args := []interface{}{tr.Tournament, tr.Rikishi, tr.Rank, tr.Wins, tr.Losses, tr.Absent}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&tr.ID, &tr.Version)
}

func (t TournamentResultModel) Get(id int64) (*TournamentResult, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, tournament, rikishi, rank, wins, losses, absent, version
		FROM tournaments_results
		WHERE id = $1`

	var tr TournamentResult

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, id).Scan(
		&tr.ID,
		&tr.Tournament,
		&tr.Rikishi,
		&tr.Rank,
		&tr.Wins,
		&tr.Losses,
		&tr.Absent,
		&tr.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &tr, nil
}

func (t TournamentResultModel) GetAll(tournament string, rank string, wins int, shikonas []string, filters Filters) ([]*TournamentResult, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, tournament, rikishi, rank, wins, losses, absent, version
		FROM tournaments_results
		WHERE (LOWER(tournament) = LOWER($1) OR $1 = '')
		AND (LOWER(rank) = LOWER($2) OR $2 = '')
		AND (rikishi = ANY($3) OR $3 = '{}')
		AND wins >= $4
		ORDER BY %s %s, id ASC
		LIMIT $5 OFFSET $6`, filters.SortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{tournament, rank, pq.Array(shikonas), wins, filters.limit(), filters.offset()}

	rows, err := t.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	tournamentsResults := []*TournamentResult{}

	for rows.Next() {
		var tournamentResult TournamentResult

		err := rows.Scan(
			&totalRecords,
			&tournamentResult.ID,
			&tournamentResult.Tournament,
			&tournamentResult.Rikishi,
			&tournamentResult.Rank,
			&tournamentResult.Wins,
			&tournamentResult.Losses,
			&tournamentResult.Absent,
			&tournamentResult.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		tournamentsResults = append(tournamentsResults, &tournamentResult)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return tournamentsResults, metadata, nil
}

func (t TournamentResultModel) Update(tr *TournamentResult) error {
	query := `
		UPDATE tournaments_results
		SET tournament = $1, rikishi = $2, rank = $3, wins = $4, losses = $5, absent = $6, version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING version`

	args := []interface{}{
		tr.Tournament,
		tr.Rikishi,
		tr.Rank,
		tr.Wins,
		tr.Losses,
		tr.Absent,
		tr.ID,
		tr.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&tr.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (t TournamentResultModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM tournaments_results WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func ValidateTournamentResult(v *validator.Validator, tr *TournamentResult, rm RikishiModel) {
	v.Check(validator.ValidTournament(tr.Tournament), "tournament", "year must be between 1900 and 2050. Month must be 3 letters. Example: 2022 Nov")

	v.Check(tr.Rikishi != "", "rikishi", "must be provided")
	v.Check(len(tr.Rikishi) <= 500, "rikishi", "must not be more than 500 bytes long")
	v.Check(rm.Exists(tr.Rikishi), "rikishi", "must exist in the database")

	v.Check(tr.Rank != "", "rank", "must be provided")
	v.Check(len(tr.Rank) <= 500, "rank", "must not be more than 500 bytes long")

	v.Check(tr.Wins >= 0 && tr.Wins <= 15, "wins", "must be between 0 and 15")
	v.Check(tr.Losses >= 0 && tr.Losses <= 15, "losses", "must be between 0 and 15")
	v.Check(tr.Absent >= 0 && tr.Absent <= 15, "absent", "must be between 0 and 15")
}
