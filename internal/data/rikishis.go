package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/corsairconstantine/sumodb/internal/validator"
	"github.com/lib/pq"
)

type Rikishi struct {
	Shikona        string   `json:"shikona"`
	HighestRank    string   `json:"highest_rank"`
	Heya           string   `json:"heya"`
	ShikonaHistory []string `json:"shikona_history"`
	Version        int32    `json:"version"`
}

type RikishiModel struct {
	DB *sql.DB
}

func (r RikishiModel) Insert(rikishi *Rikishi) error {
	query := `
		INSERT INTO rikishis (shikona, highest_rank, heya, shikona_history)
		VALUES ($1, $2, $3, $4)
		RETURNING version`

	args := []interface{}{rikishi.Shikona, rikishi.HighestRank, rikishi.Heya, pq.Array(rikishi.ShikonaHistory)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&rikishi.Version)
}

func (r RikishiModel) Get(shikona string) (*Rikishi, error) {
	if shikona == "" {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT shikona, highest_rank, heya, shikona_history, version
		FROM rikishis
		WHERE shikona = $1`

	var rikishi Rikishi

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, shikona).Scan(
		&rikishi.Shikona,
		&rikishi.HighestRank,
		&rikishi.Heya,
		pq.Array(&rikishi.ShikonaHistory),
		&rikishi.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &rikishi, nil
}

func (r RikishiModel) Update(rikishi *Rikishi) error {
	var oldShikona string = rikishi.ShikonaHistory[len(rikishi.ShikonaHistory)-1]
	if rikishi.Shikona != oldShikona {
		rikishi.ShikonaHistory = append(rikishi.ShikonaHistory, rikishi.Shikona)
	}
	query := `
		UPDATE rikishis
		SET shikona = $1, highest_rank = $2, heya = $3, shikona_history = $4, version = version + 1
		WHERE shikona = $5 AND version = $6
		RETURNING version`

	args := []interface{}{
		rikishi.Shikona,
		rikishi.HighestRank,
		rikishi.Heya,
		pq.Array(rikishi.ShikonaHistory),
		oldShikona,
		rikishi.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(&rikishi.Version)
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

func (r RikishiModel) Delete(shikona string) error {
	if shikona == "" {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM rikishis
		WHERE shikona = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, query, shikona)
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

func (r RikishiModel) Exists(shikona string) bool {
	var exists bool
	query := `SELECT exists (SELECT true FROM rikishis WHERE shikona = $1)`
	r.DB.QueryRow(query, shikona).Scan(&exists)

	return exists
}

func ValidateRikishi(v *validator.Validator, rikishi *Rikishi) {
	v.Check(rikishi.Shikona != "", "shikona", "must be provided")
	v.Check(len(rikishi.Shikona) <= 500, "shikona", "must not be more than 500 bytes long")

	v.Check(rikishi.HighestRank != "", "highest rank", "must be provided")
	v.Check(len(rikishi.HighestRank) <= 500, "rank", "must not be more than 500 bytes long")

	v.Check(rikishi.Heya != "", "heya", "must be provided")
	v.Check(len(rikishi.Heya) <= 500, "rank", "must not be more than 500 bytes long")

	v.Check(rikishi.ShikonaHistory != nil, "shikona history", "must be provided")
	v.Check(len(rikishi.ShikonaHistory) >= 1, "shikona history", "must contain at least 1 shikona")
	v.Check(validator.Unique(rikishi.ShikonaHistory), "shikona history", "must not contain duplicate values")
}
