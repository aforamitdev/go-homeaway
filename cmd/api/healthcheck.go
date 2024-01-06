package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelop{
		"status": "available",
		"system_info": map[string]string{
			"env":     app.config.evn,
			"version": version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The Server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}
