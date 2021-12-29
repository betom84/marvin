package hue_test

import (
	"marvin/config"
	"marvin/devices/hue"
	"marvin/metrics"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	promtest "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestLight(t *testing.T) {
	tt := []struct {
		name    string
		light   hue.Light
		metrics map[string]int
	}{
		{
			name:  "Change state of light 3",
			light: hue.NewLight(3, config.Get().PhilipsHueHost, config.Get().PhilipsHueUser),
			metrics: map[string]int{
				"application_device_operation_duration_milliseconds": 2,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			prom := prometheus.NewRegistry()
			metrics.Instance().Register(prom)

			originState, err := tc.light.State()
			assert.NoError(t, err)

			changedState, err := tc.light.SetState(!originState)
			assert.NoError(t, err)
			assert.NotEqual(t, originState, changedState)

			time.Sleep(2 * time.Second)

			changedState, err = tc.light.SetState(originState)
			assert.NoError(t, err)
			assert.Equal(t, originState, changedState)

			for em, ec := range tc.metrics {
				c, err := promtest.GatherAndCount(prom, em)
				assert.NoError(t, err)
				assert.Equal(t, ec, c)
			}
		})
	}
}
