package api

import (
	"io/ioutil"
	"net/http"
)

func HandleEndpointsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("config/endpoints.json")
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(content)
	}
}
