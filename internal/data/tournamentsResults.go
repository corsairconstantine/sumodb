package data

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/corsairconstantine/sumodb/internal/validator"
)

type TournamentResult struct {
	ID         int64  `json:"id"`
	Tournament Date   `json:"tournament"`
	Rikishi    string `json:"rikishi"`
	Rank       string `json:"rank"`
	Result     string `json:"result"`
	Version    int32  `json:"version"`
}

type TournamentResultModel struct {
	DB *sql.DB
}

func (t TournamentResultModel) Insert(tr *TournamentResult) error {
	query := `
		INSERT INTO tournaments_results (tournament, rikishi, rank, result)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version`

	args := []interface{}{tr.Tournament.Time, tr.Rikishi, tr.Rank, tr.Result}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&tr.ID, &tr.Version)
}

func (t TournamentResultModel) Get(id int64) (*TournamentResult, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, tournament, rikishi, rank, result, version
		FROM tournaments_results
		WHERE id = $1`

	var tr TournamentResult

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, id).Scan(
		&tr.ID,
		&tr.Tournament.Time,
		&tr.Rikishi,
		&tr.Rank,
		&tr.Result,
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

func (t TournamentResultModel) Update(tr *TournamentResult) error {
	query := `
		UPDATE tournaments_results
		SET tournament = $1, rikishi = $2, rank = $3, result = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []interface{}{
		tr.Tournament.Time,
		tr.Rikishi,
		tr.Rank,
		tr.Result,
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
	v.Check(!tr.Tournament.Before(time.Date(1900, 0, 0, 0, 0, 0, 0, time.UTC)), "tournament", "date must be after year 1900")
	v.Check(!tr.Tournament.After(time.Date(2050, 0, 0, 0, 0, 0, 0, time.UTC)), "tournament", "date must be before year 2050")

	v.Check(tr.Rikishi != "", "rikishi", "must be provided")
	v.Check(len(tr.Rikishi) <= 500, "rikishi", "must not be more than 500 bytes long")
	v.Check(rm.Exists(tr.Rikishi), "rikishi", "must exist in the database")

	v.Check(tr.Rank != "", "rank", "must be provided")
	v.Check(len(tr.Rank) <= 500, "rank", "must not be more than 500 bytes long")

	v.Check(tr.Result != "", "result", "must be provided")
	v.Check(len(tr.Result) <= 500, "result", "must not be more than 500 bytes long")

	var areInts, isLessThan15 bool = true, true
	scores := strings.Split(tr.Result, "-")
	var sum int
	for _, v := range scores {
		i, err := strconv.Atoi(v)
		if err != nil {
			areInts = false
		}
		sum += i
	}
	if sum > 15 {
		isLessThan15 = false
	}

	v.Check(areInts, "result", "must be integers separated by '-'")
	v.Check(isLessThan15, "result", "sum of wins and losses must be less than 15")
}
