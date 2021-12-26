package api

import (
	"encoding/json"
	"marvin/alexa"
	"net/http"
)

func HandleAlexaStateGet(a *alexa.Server) http.HandlerFunc {
	var response = struct {
		State string `json:"state"`
	}{}

	return func(w http.ResponseWriter, r *http.Request) {
		response.State = "stopped"
		if a.IsRunning() {
			response.State = "running"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func HandleAlexaStateSet(a *alexa.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		switch r.URL.Query().Get("set") {
		case "running":
			err = a.Start()
		case "stopped":
			err = a.Stop()
		}

		if err != nil {
			panic(err)
		}
	}
}
