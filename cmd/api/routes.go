package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healshcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/rikishis", app.listRikishisHandler)
	router.HandlerFunc(http.MethodPost, "/v1/rikishis", app.createRikishiHandler)
	router.HandlerFunc(http.MethodGet, "/v1/rikishis/:shikona", app.showRikishiHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/rikishis/:shikona", app.updateRikishiHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/rikishis/:shikona", app.deleteRikishiHandler)

	router.HandlerFunc(http.MethodGet, "/v1/tournamentsresults", app.listTournamentsResultsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tournamentsresults", app.createTournamentResultHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tournamentsresults/:id", app.showTournamentResultHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/tournamentsresults/:id", app.updateTournamentResultHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/tournamentsresults/:id", app.deleteTournamentResultHandler)

	router.HandlerFunc(http.MethodGet, "/v1/bouts", app.listBoutsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/bouts", app.createBoutHandler)
	router.HandlerFunc(http.MethodGet, "/v1/bouts/:id", app.showBoutHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/bouts/:id", app.updateBoutHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/bouts/:id", app.deleteBoutHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return app.recoverPanic(app.rateLimit(router))
}
