package hue_test

import (
	"marvin/config"
	"marvin/devices/hue"
	"testing"
	"time"
)

// todo
func _TestLight(t *testing.T) {
	tt := []struct {
		name  string
		light hue.Light
	}{
		{
			name:  "Change state of light 3",
			light: hue.NewLight(3, config.Get().PhilipsHueHost, config.Get().PhilipsHueUser),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			originState, err := tc.light.State()
			if err != nil {
				t.Fatalf("got an error; %v", err)
			}

			changedState, err := tc.light.SetState(!originState)
			if err != nil {
				t.Fatalf("got an error; %v", err)
			}

			if changedState == originState {
				t.Fatalf("state was not changed to %v", !originState)
			}

			time.Sleep(1 * time.Second)

			changedState, err = tc.light.SetState(originState)
			if err != nil {
				t.Fatalf("got an error; %v", err)
			}

			if changedState != originState {
				t.Fatalf("state was not changed back to origin %v", originState)
			}
		})
	}
}
