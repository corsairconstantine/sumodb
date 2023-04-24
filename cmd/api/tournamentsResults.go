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
		Tournament data.Date `json:"tournament"`
		Rikishi    string    `json:"rikishi"`
		Rank       string    `json:"rank"`
		Result     string    `json:"result"`
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
		Result:     input.Result,
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
	//possible location "/v1/rikishis/shikona/tournamentsresults/id"

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
		Tournament *data.Date `json:"tournament"`
		Rikishi    *string    `json:"rikishi"`
		Rank       *string    `json:"rank"`
		Result     *string    `json:"result"`
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

	if input.Result != nil {
		tr.Result = *input.Result
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
