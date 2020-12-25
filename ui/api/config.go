package api

import (
	"encoding/json"
	"marvin/config"
	"net/http"
)

type configGetResponse struct {
	Current config.Configuration `json:"current"`
	Raw     config.Configuration `json:"raw"`
	Default config.Configuration `json:"default"`
}

func HandleConfigGet() func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		configGetResponse := configGetResponse{
			Current: config.Get(),
			Raw:     config.GetRaw(),
			Default: config.GetDefault(),
		}

		content, err := json.MarshalIndent(configGetResponse, "", "   ")
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(content)
		return nil
	}
}
