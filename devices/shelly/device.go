package shelly

import (
	"encoding/json"
	"fmt"
	"marvin/metrics"
	"net/http"
	"strings"
)

type Device interface {
	State() (bool, error)
	SetState(bool) (bool, error)
}

type ShellyOne struct {
	id string
}

func (d ShellyOne) State() (bool, error) {
	var err error
	defer metrics.CollectDeviceOperationDuration(d.id, "State", err)()

	r, err := http.Get(fmt.Sprintf("http://%s/relay/0", d.id))
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	resp := struct {
		IsON bool `json:"ison"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return false, err
	}

	return resp.IsON, nil
}

func (d ShellyOne) SetState(value bool) (bool, error) {
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

	resp := struct {
		IsON bool `json:"ison"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return false, err
	}

	return resp.IsON, nil
}

type ShellyPlus2PM struct {
	id string
}

func (d ShellyPlus2PM) State() (bool, error) {
	var err error
	defer metrics.CollectDeviceOperationDuration(d.id, "State", err)()

	r, err := http.Get(fmt.Sprintf("http://%s/rpc/Cover.GetStatus?id=0", d.id))
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	resp := struct {
		State string `json:"state"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return false, err
	}

	return resp.State == "open", nil
}

func (d ShellyPlus2PM) SetState(value bool) (bool, error) {
	var err error
	defer metrics.CollectDeviceOperationDuration(d.id, "SetState", err)()

	op := "Close"
	if value {
		op = "Open"
	}

	r, err := http.Get(fmt.Sprintf("http://%s/rpc/Cover.%s?id=0", d.id, op))
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	return value, nil
}

func NewDevice(id string) (Device, error) {
	var d Device
	var err error

	switch strings.Split(id, "-")[0] {
	case "shelly1":
		d = ShellyOne{id}
	case "shellyplus2pm":
		d = ShellyPlus2PM{id}
	default:
		err = fmt.Errorf("shelly device type not supported")
	}

	return d, err
}
