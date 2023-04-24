package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/corsairconstantine/sumodb/internal/data"
	"github.com/corsairconstantine/sumodb/internal/validator"
)

func (app *application) createRikishiHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Shikona        string   `json:"shikona"`
		HighestRank    string   `json:"highest_rank"`
		Heya           string   `json:"heya"`
		ShikonaHistory []string `json:"shikona_history"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	rikishi := &data.Rikishi{
		Shikona:        input.Shikona,
		HighestRank:    input.HighestRank,
		Heya:           input.Heya,
		ShikonaHistory: input.ShikonaHistory,
	}

	v := validator.New()

	if data.ValidateRikishi(v, rikishi); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Rikishis.Insert(rikishi)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/rikishis/%s", strings.ReplaceAll(rikishi.Shikona, " ", "-")))

	err = app.writeJSON(w, http.StatusCreated, envelope{"rikishi": rikishi}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showRikishiHandler(w http.ResponseWriter, r *http.Request) {
	shikona, err := app.readShikonaParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	rikishi, err := app.models.Rikishis.Get(shikona)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"rikishi": rikishi}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateRikishiHandler(w http.ResponseWriter, r *http.Request) {
	shikona, err := app.readShikonaParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	rikishi, err := app.models.Rikishis.Get(shikona)
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
		Shikona        *string  `json:"shikona"`
		NewShikona     *string  `json:"new_shikona"`
		HighestRank    *string  `json:"highest_rank"`
		Heya           *string  `json:"heya"`
		ShikonaHistory []string `json:"shikona_history"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Shikona != nil {
		rikishi.Shikona = *input.Shikona
	}

	if input.NewShikona != nil {
		rikishi.Shikona = *input.NewShikona
	}

	if input.HighestRank != nil {
		rikishi.HighestRank = *input.HighestRank
	}

	if input.Heya != nil {
		rikishi.Heya = *input.Heya
	}

	if input.ShikonaHistory != nil {
		rikishi.ShikonaHistory = input.ShikonaHistory
	}

	v := validator.New()

	if data.ValidateRikishi(v, rikishi); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Rikishis.Update(rikishi)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"rikishi": rikishi}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteRikishiHandler(w http.ResponseWriter, r *http.Request) {
	shikona, err := app.readShikonaParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Rikishis.Delete(shikona)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "rikishi successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

/*
func (app *application) listRikishiHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Shikona string
		ShikonaHistory []string
		Page int
		PageSize int
		Sort string
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Shikona
}*/
