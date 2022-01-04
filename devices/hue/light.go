package hue

import (
	"encoding/json"
	"fmt"
	"marvin/metrics"
	"net/http"
	"strings"
)

// Light represents a hue light based on his lightID
type Light struct {
	lightID int
	host    string
	user    string
}

type getStateResponse struct {
	State struct {
		On        bool      `json:"on"`
		Bri       int       `json:"bri"`
		Hue       int       `json:"hue"`
		Sat       int       `json:"sat"`
		Effect    string    `json:"effect"`
		Xy        []float64 `json:"xy"`
		Ct        int       `json:"ct"`
		Alert     string    `json:"alert"`
		Colormode string    `json:"colormode"`
		Reachable bool      `json:"reachable"`
	} `json:"state"`
	Swupdate struct {
		State       string      `json:"state"`
		Lastinstall interface{} `json:"lastinstall"`
	} `json:"swupdate"`
	Type             string `json:"type"`
	Name             string `json:"name"`
	Modelid          string `json:"modelid"`
	Manufacturername string `json:"manufacturername"`
	Uniqueid         string `json:"uniqueid"`
	Swversion        string `json:"swversion"`
	Swconfigid       string `json:"swconfigid"`
	Productid        string `json:"productid"`
}

// State retrieves the current state of the hue light
func (l Light) State() (bool, error) {
	var err error
	defer metrics.CollectDeviceOperationDuration(fmt.Sprintf("hue-%d", l.lightID), "State", err)()

	r, err := http.Get(l.getURL())
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	var resp getStateResponse
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return false, err
	}

	return resp.State.On, nil
}

func (l Light) getURL() string {
	return fmt.Sprintf("http://%s/api/%s/lights/%d", l.host, l.user, l.lightID)
}

type setStateResponse []struct {
	Success map[string]interface{} `json:"success"`
}

// SetState change the current state of the light
func (l Light) SetState(value bool) (bool, error) {
	var err error
	defer metrics.CollectDeviceOperationDuration(fmt.Sprintf("hue-%d", l.lightID), "SetState", err)()

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/state", l.getURL()),
		strings.NewReader(fmt.Sprintf(`{ "on": %v }`, value)))

	if err != nil {
		return false, err
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	var resp setStateResponse
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return false, err
	}

	return resp[0].Success[fmt.Sprintf("/lights/%d/state/on", l.lightID)].(bool), nil
}

// NewLight creates a hue light instance based on the lightID
func NewLight(lightID int, host string, user string) Light {
	return Light{lightID: lightID, host: host, user: user}
}
