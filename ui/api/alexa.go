package api

import (
	"encoding/json"
	"marvin/alexa"
	"net/http"
)

func HandleAlexaStateGet(a *alexa.Server) func(w http.ResponseWriter, r *http.Request) error {
	var response = struct {
		State string `json:"state"`
	}{}

	return func(w http.ResponseWriter, r *http.Request) error {
		response.State = "stopped"
		if a.IsRunning() {
			response.State = "running"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		return nil
	}
}

func HandleAlexaStateSet(a *alexa.Server) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error

		switch r.URL.Query().Get("set") {
		case "running":
			err = a.Start()
		case "stopped":
			err = a.Stop()
		}

		return err
	}
}
