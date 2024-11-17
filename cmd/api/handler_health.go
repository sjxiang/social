package main

import (
	"net/http"

	"github.com/sjxiang/social/internal/store"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": app.config.version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) createDogWithBuilder(w http.ResponseWriter, r *http.Request) {
	
	p, err := store.NewPetBuilder().
		SetSpecies("狗").
		SetBreed("中华田园犬").
		SetWeight(15).
		SetDescription("性情温顺").
		SetColor("黑").
		SetAge(3).
		SetAgeEstimated(true).
		Build()

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	
	if err := app.jsonResponse(w, http.StatusOK, p); err != nil {
		app.internalServerError(w, r, err)
	}
}
