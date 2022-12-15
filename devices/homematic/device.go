package homematic

import (
	"encoding/xml"
	"fmt"
	"io"
	"marvin/metrics"
	"net/http"
	"strconv"

	"golang.org/x/text/encoding/charmap"
)

// A Device represents an homematic device represented by the datapoint ise_id
type Device struct {
	iseID int
	host  string
}

// State retrieves the current state of the homematic device
func (d Device) State() (bool, error) {
	var err error
	defer metrics.CollectDeviceOperationDuration(fmt.Sprintf("homematic-%d", d.iseID), "State", err)()

	r, err := http.Get(fmt.Sprintf("http://%s/addons/xmlapi/state.cgi?datapoint_id=%d", d.host, d.iseID))
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	var s struct {
		XMLName   xml.Name `xml:"state"`
		Datapoint struct {
			XMLName xml.Name `xml:"datapoint"`
			IseID   int      `xml:"ise_id,attr"`
			Value   string   `xml:"value,attr"`
		}
	}

	err = d.parseResponse(r.Body, &s)
	if err != nil {
		return false, err
	}

	return d.stateToBool(s.Datapoint.Value), nil
}

// SetState change the current state of the homematic device
func (d Device) SetState(value bool) (bool, error) {
	var err error
	defer metrics.CollectDeviceOperationDuration(fmt.Sprintf("homematic-%d", d.iseID), "SetState", err)()

	r, err := http.Get(fmt.Sprintf("http://%s/addons/xmlapi/statechange.cgi?ise_id=%d&new_value=%s", d.host, d.iseID, d.boolToState(value)))
	if err != nil {
		return false, err
	}
	defer r.Body.Close()

	var res struct {
		XMLName xml.Name `xml:"result"`
		Changed struct {
			XMLName  xml.Name `xml:"changed"`
			ID       int      `xml:"id,attr"`
			NewValue string   `xml:"new_value,attr"`
		}
	}

	err = d.parseResponse(r.Body, &res)
	if err != nil {
		return false, err
	}

	return d.stateToBool(res.Changed.NewValue), nil
}

func (d Device) boolToState(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func (d Device) stateToBool(v string) bool {
	if b, err := strconv.ParseBool(v); err == nil {
		return b
	}

	if f, err := strconv.ParseFloat(v, 32); err == nil {
		return f != 0.0
	}

	return false
}

func (d Device) parseResponse(body io.Reader, target interface{}) (err error) {
	decoder := xml.NewDecoder(body)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return charmap.ISO8859_1.NewDecoder().Reader(input), nil
	}
	err = decoder.Decode(target)
	return
}

// NewDevice creates a homematic device instance based on the given id (iseID)
func NewDevice(id string, host string) (Device, error) {
	iseID, err := strconv.Atoi(id)
	if err != nil {
		return Device{}, fmt.Errorf("unable to instantiate device; %s", err)
	}

	return Device{iseID: iseID, host: host}, nil
}
