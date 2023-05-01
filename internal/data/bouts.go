package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/corsairconstantine/sumodb/internal/validator"
	"github.com/lib/pq"
)

type Bout struct {
	ID         int64
	Tournament string
	Day        string
	Winner     string
	Loser      string
	Kimarite   string
	Version    int32
}

type BoutModel struct {
	DB *sql.DB
}

func (b BoutModel) Insert(bout *Bout) error {
	query := `
		INSERT INTO bouts (tournament, day, winner, loser, kimarite)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version`

	args := []interface{}{bout.Tournament, bout.Day, bout.Winner, bout.Loser, bout.Kimarite}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return b.DB.QueryRowContext(ctx, query, args...).Scan(&bout.ID, &bout.Version)
}

func (b BoutModel) Get(id int64) (*Bout, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, tournament, day, winner, loser, kimarite, version
		FROM bouts
		WHERE id = $1`

	var bout Bout

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, id).Scan(
		&bout.ID,
		&bout.Tournament,
		&bout.Day,
		&bout.Winner,
		&bout.Loser,
		&bout.Kimarite,
		&bout.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &bout, nil
}

func (b BoutModel) GetAll(tournament, day, kimarite string, rikishi1, rikishi2 []string) ([]*Bout, error) {
	query := `
		SELECT id, tournament, day, winner, loser, kimarite, version
		FROM bouts
		WHERE (LOWER(tournament) = LOWER($1) OR $1 = '')
		AND (day = $2 OR $2 = '')
		AND (kimarite = $3 OR $3 = '')
		AND (winner = ANY($4) OR loser = ANY($4) OR $4 = '{}')
		AND (winner = ANY($5) OR loser = ANY($5) OR $5 = '{}')`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{tournament, day, kimarite, pq.Array(rikishi1), pq.Array(rikishi2)}

	rows, err := b.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bouts := []*Bout{}

	for rows.Next() {
		var bout Bout

		err := rows.Scan(
			&bout.ID,
			&bout.Tournament,
			&bout.Day,
			&bout.Winner,
			&bout.Loser,
			&bout.Kimarite,
			&bout.Version,
		)

		if err != nil {
			return nil, err
		}

		bouts = append(bouts, &bout)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bouts, nil
}

func (b BoutModel) Update(bout *Bout) error {
	query := `
		UPDATE bouts
		SET tournament = $1, day = $2, winner = $3, loser = $4, kimarite = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version`

	args := []interface{}{
		bout.Tournament,
		bout.Day,
		bout.Winner,
		bout.Loser,
		bout.Kimarite,
		bout.ID,
		bout.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, args...).Scan(&bout.Version)
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

func (b BoutModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM bouts WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := b.DB.ExecContext(ctx, query, id)
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

func ValidateBout(v *validator.Validator, b *Bout, rm RikishiModel) {
	v.Check(validator.ValidTournament(b.Tournament), "tournament", "year must be between 1900 and 2050. Month must be 3 letters. Example: 2022 Nov")

	v.Check(validator.ValidDay(b.Day), "day", "must be a number from 1 to 15. Alternatively can be 'Playoff'")

	v.Check(b.Winner != "", "winner", "must be provided")
	v.Check(len(b.Winner) <= 500, "winner", "must not be more than 500 bytes long")
	v.Check(rm.Exists(b.Winner), "winner", "must exist in the database")

	v.Check(b.Loser != "", "loser", "must be provided")
	v.Check(len(b.Loser) <= 500, "loser", "must not be more than 500 bytes long")
	v.Check(rm.Exists(b.Loser), "loser", "must exist in the database")

	v.Check(len(b.Kimarite) <= 500, "kimarite", "must not be more than 500 bytes long")
}
