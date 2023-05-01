package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/corsairconstantine/sumodb/internal/data"
	"github.com/corsairconstantine/sumodb/internal/validator"
)

func (app *application) createBoutHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Tournament string `json:"tournament"`
		Day        string `json:"day"`
		Winner     string `json:"winner"`
		Loser      string `json:"loser"`
		Kimarite   string `json:"kimarite"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	bout := &data.Bout{
		Tournament: input.Tournament,
		Day:        input.Day,
		Winner:     input.Winner,
		Loser:      input.Loser,
		Kimarite:   input.Kimarite,
	}

	v := validator.New()
	if data.ValidateBout(v, bout, app.models.Rikishis); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Bouts.Insert(bout)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/bouts/%d", bout.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"bout": bout}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showBoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	bout, err := app.models.Bouts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"bout": bout}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateBoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	bout, err := app.models.Bouts.Get(id)
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
		Day        *string `json:"day"`
		Winner     *string `json:"winner"`
		Loser      *string `json:"loser"`
		Kimarite   *string `json:"kimarite"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Tournament != nil {
		bout.Tournament = *input.Tournament
	}

	if input.Day != nil {
		bout.Day = *input.Day
	}

	if input.Winner != nil {
		bout.Winner = *input.Winner
	}

	if input.Loser != nil {
		bout.Loser = *input.Loser
	}

	if input.Kimarite != nil {
		bout.Kimarite = *input.Kimarite
	}

	v := validator.New()
	if data.ValidateBout(v, bout, app.models.Rikishis); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Bouts.Update(bout)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"bout": bout}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Bouts.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "bout successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listBoutsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Tournament string
		Day        string
		Rikishi1   string
		Rikishi2   string
		Kimarite   string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Tournament = app.readString(qs, "tournament", "")
	input.Day = app.readString(qs, "day", "")
	input.Kimarite = app.readString(qs, "kimarite", "")
	input.Rikishi1 = app.readString(qs, "rikishi1", "")
	input.Rikishi1 = app.readString(qs, "rikishi2", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "tournament", "-tournament"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	r1, err := app.models.Rikishis.GetShikonaHistory(input.Rikishi1)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	r2, err := app.models.Rikishis.GetShikonaHistory(input.Rikishi2)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	bouts, err := app.models.Bouts.GetAll(input.Tournament, input.Day, input.Kimarite, r1, r2)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"bouts": bouts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
