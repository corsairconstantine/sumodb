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

func (r RikishiModel) GetAll(shikona, highestRank, heya string, filters Filters) ([]*Rikishi, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), shikona, highest_rank, heya, shikona_history, version
		FROM rikishis
		WHERE (array_to_string(shikona_history, ',') @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (LOWER(highest_rank) = LOWER($2) OR $2 = '')
		AND (LOWER(heya) = LOWER($3) OR $3 = '')
		ORDER BY %s %s, shikona ASC
		LIMIT $4 OFFSET $5`, filters.SortColumn(), filters.sortDirection())

	//need a custom sort by rank: yokozuna->ozeki->sekiwake->komusubi->maegashira->
	//juryo->makushita->sandanme->jonidan->jonokuchi->mae-zumo

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{shikona, highestRank, heya, filters.limit(), filters.offset()}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	rikishis := []*Rikishi{}

	for rows.Next() {
		var rikishi Rikishi

		err := rows.Scan(
			&totalRecords,
			&rikishi.Shikona,
			&rikishi.HighestRank,
			&rikishi.Heya,
			pq.Array(&rikishi.ShikonaHistory),
			&rikishi.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		rikishis = append(rikishis, &rikishi)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return rikishis, metadata, nil
}

func (r RikishiModel) GetShikonaHistory(shikona string) ([]string, error) {
	if shikona == "" {
		return []string{}, nil
	}

	query := `
		SELECT shikona_history
		FROM rikishis
		WHERE array_to_string(shikona_history, ',') @@ plainto_tsquery('simple', $1)`

	rows, err := r.DB.Query(query, shikona)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	var shikonas []string

	for rows.Next() {
		var shikona_history []string
		err := rows.Scan(pq.Array(&shikona_history))
		if err != nil {
			return nil, err
		}
		shikonas = append(shikonas, shikona_history...)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shikonas, nil
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
