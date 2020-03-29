package api

import (
	"io/ioutil"
	"net/http"
)

func HandleConfigGet() func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		content, err := ioutil.ReadFile("config/config.json")
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(content)
		return nil
	}
}
