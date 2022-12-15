package shelly

import (
	"encoding/json"
	"fmt"
	"marvin/metrics"
	"net/http"
)

type Device struct {
	id string
}

type APIResponse struct {
	IsON bool `json:"ison"`
}

func (d Device) State() (bool, error) {
	var err error
	defer metrics.CollectDeviceOperationDuration(d.id, "State", err)()

	r, err := http.Get(fmt.Sprintf("http://%s/relay/0", d.id))
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	resp := APIResponse{}
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return false, err
	}

	return resp.IsON, nil
}

func (d Device) SetState(value bool) (bool, error) {
	var err error
	defer metrics.CollectDeviceOperationDuration(d.id, "SetState", err)()

	v := "off"
	if value {
		v = "on"
	}

	r, err := http.Get(fmt.Sprintf("http://%s/relay/0?turn=%s", d.id, v))
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	resp := APIResponse{}
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return false, err
	}

	return resp.IsON, nil
}

func NewDevice(id string) (Device, error) {
	return Device{id}, nil
}
