package api_test

import (
	"marvin/ui/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	err := api.HandleConfigGet()(w, r)
	assert.NoError(t, err)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Greater(t, w.Body.Len(), 0)

	t.Log(w.Body.ReadString(0x0))
}
