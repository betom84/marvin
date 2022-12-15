package shelly_test

import (
	"marvin/devices/shelly"
	"marvin/metrics"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	promtest "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDevice(t *testing.T) {
	device, _ := shelly.NewDevice("shelly1-244cab43abbe")

	tt := []struct {
		name   string
		device shelly.Device
	}{
		{
			name:   "Change state of shelly device",
			device: device,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			prom := prometheus.NewRegistry()
			metrics.Instance().Register(prom)

			originState, err := tc.device.State()
			assert.NoError(t, err)

			changedState, err := tc.device.SetState(!originState)
			assert.NoError(t, err)
			assert.NotEqual(t, originState, changedState)

			time.Sleep(2 * time.Second)

			changedState, err = tc.device.SetState(originState)
			assert.NoError(t, err)
			assert.Equal(t, originState, changedState)

			c, err := promtest.GatherAndCount(prom, "application_device_operation_duration_milliseconds")
			assert.NoError(t, err)
			assert.Equal(t, 2, c)
		})
	}
}
