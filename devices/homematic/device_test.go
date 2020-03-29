package homematic_test

import (
	"marvin/config"
	"marvin/devices/homematic"
	"testing"
	"time"
)

func TestDevice(t *testing.T) {
	tt := []struct {
		name   string
		device homematic.Device
	}{
		{
			name:   "Change state of datapointIseID",
			device: homematic.NewDevice(2195, config.Get().HomematicHost),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			originState, err := tc.device.State()
			if err != nil {
				t.Fatalf("got an error; %v", err)
			}

			changedState, err := tc.device.SetState(!originState)
			if err != nil {
				t.Fatalf("got an error; %v", err)
			}

			if changedState == originState {
				t.Fatalf("state was not changed to %v", !originState)
			}

			time.Sleep(2 * time.Second)

			changedState, err = tc.device.SetState(originState)
			if err != nil {
				t.Fatalf("got an error; %v", err)
			}

			if changedState != originState {
				t.Fatalf("state was not changed back to origin %v", originState)
			}
		})
	}
}
