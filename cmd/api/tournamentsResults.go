package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/corsairconstantine/sumodb/internal/data"
	"github.com/corsairconstantine/sumodb/internal/validator"
)

func (app *application) createTournamentResultHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Tournament string `json:"tournament"`
		Rikishi    string `json:"rikishi"`
		Rank       string `json:"rank"`
		Wins       int32  `json:"wins"`
		Losses     int32  `json:"losses"`
		Absent     int32  `json:"absent"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tr := &data.TournamentResult{
		Tournament: input.Tournament,
		Rikishi:    input.Rikishi,
		Rank:       input.Rank,
		Wins:       input.Wins,
		Losses:     input.Losses,
		Absent:     input.Absent,
	}

	v := validator.New()
	if data.ValidateTournamentResult(v, tr, app.models.Rikishis); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.TournamentsResults.Insert(tr)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tournamentsresults/%d", tr.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"tournament_result": tr}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showTournamentResultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tr, err := app.models.TournamentsResults.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tournament_result": tr}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTournamentResultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tr, err := app.models.TournamentsResults.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Tournament *string `json:"tournament"`
		Rikishi    *string `json:"rikishi"`
		Rank       *string `json:"rank"`
		Wins       *int32  `json:"result"`
		Losses     *int32  `json:"losses"`
		Absent     *int32  `json:"absent"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Tournament != nil {
		tr.Tournament = *input.Tournament
	}

	if input.Rikishi != nil {
		tr.Rikishi = *input.Rikishi
	}

	if input.Rank != nil {
		tr.Rank = *input.Rank
	}

	if input.Wins != nil {
		tr.Wins = *input.Wins
	}

	if input.Losses != nil {
		tr.Losses = *input.Losses
	}

	if input.Absent != nil {
		tr.Absent = *input.Absent
	}

	v := validator.New()
	if data.ValidateTournamentResult(v, tr, app.models.Rikishis); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.TournamentsResults.Update(tr)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tournament_result": tr}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTournamentResultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.TournamentsResults.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "tournament record successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listTournamentsResultsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Tournament string
		Rikishi    string
		Rank       string
		Wins       int
		Loses      int
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Tournament = app.readString(qs, "tournament", "")
	input.Rikishi = app.readString(qs, "rikishi", "")
	input.Rank = app.readString(qs, "rank", "")
	input.Wins = app.readInt(qs, "wins", 0, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "tournament", "rikishi", "rank", "wins", "-id", "-tournament", "-rikishi", "-rank", "-wins"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	shikonas, err := app.models.Rikishis.GetShikonaHistory(input.Rikishi)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	trs, metadata, err := app.models.TournamentsResults.GetAll(input.Tournament, input.Rank, input.Wins, shikonas, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tournaments_results": trs, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
