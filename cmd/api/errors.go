package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		url    = r.URL.RequestURI()
	)
	app.logger.Error(err.Error(), method, url)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, message any, status int) {
	errMsg := envelope{
		"error": message,
	}
	err := app.writeJSON(w, errMsg, status, nil)
	if err != nil {
		app.logger.Error(err.Error())
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error())
	msg := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, msg, http.StatusInternalServerError)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "the requested resource could not be found"
	app.errorResponse(w, r, msg, http.StatusNotFound)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, msg, http.StatusMethodNotAllowed)
}
