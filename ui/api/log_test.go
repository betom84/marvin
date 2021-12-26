package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"marvin/config"
	"marvin/logger"
	"marvin/ui/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestLogEntries(t *testing.T) {
	prepareLogfile()

	type line struct {
		line   int
		suffix string
	}

	type response struct {
		Lines []string `json:"lines"`
	}

	tt := []struct {
		params            string
		expectedLineCount int
		expectedLines     []line
	}{
		{
			params:            "",
			expectedLineCount: 10,
			expectedLines: []line{
				line{
					line:   0,
					suffix: "This is log entry 91",
				},
				line{
					line:   9,
					suffix: "This is log entry 100",
				},
			},
		},
		{
			params:            "?limit=2",
			expectedLineCount: 2,
			expectedLines: []line{
				line{
					line:   0,
					suffix: "This is log entry 99",
				},
				line{
					line:   1,
					suffix: "This is log entry 100",
				},
			},
		},
		{
			params:            "?limit=2&offset=1",
			expectedLineCount: 2,
			expectedLines: []line{
				line{
					line:   0,
					suffix: "This is log entry 98",
				},
				line{
					line:   1,
					suffix: "This is log entry 99",
				},
			},
		},
		{
			params:            "?limit=2&offset=98",
			expectedLineCount: 2,
			expectedLines: []line{
				line{
					line:   0,
					suffix: "This is log entry 1",
				},
				line{
					line:   1,
					suffix: "This is log entry 2",
				},
			},
		},
		{
			params:            "?limit=2&offset=99",
			expectedLineCount: 1,
			expectedLines: []line{
				line{
					line:   0,
					suffix: "This is log entry 1",
				},
			},
		},
		{
			params:            "?limit=2&offset=100",
			expectedLineCount: 0,
			expectedLines:     []line{},
		},
		{
			params:            "?limit=200",
			expectedLineCount: 100,
			expectedLines:     []line{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.params, func(t *testing.T) {
			var response = response{}
			invokeHttp(t, tc.params, &response)
			assert.Len(t, response.Lines, tc.expectedLineCount)

			for _, el := range tc.expectedLines {
				l := response.Lines[el.line]
				assert.Truef(t, strings.HasSuffix(l, el.suffix), "'%s' was unexpected", l)
			}
		})
	}
}

func prepareLogfile() {
	sampleLog, _ := ioutil.TempFile("", "sampleLog")
	config.Set(config.Configuration{
		Log: sampleLog.Name(),
	})

	startTime := time.Now()
	for entry := 1; entry <= 100; entry++ {
		io.WriteString(sampleLog, fmt.Sprintf(
			"%s This is log entry %d\n",
			startTime.Add(time.Duration(entry)*time.Second).Format(time.RFC822),
			entry,
		))
	}
}

func invokeHttp(t *testing.T, query string, response interface{}) int {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", fmt.Sprintf("/test%s", query), nil)
	api.HandleLogGet()(w, r)

	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	err := json.NewDecoder(w.Body).Decode(response)
	assert.NoError(t, err)

	return w.Code
}

func TestHandleLogSocket(t *testing.T) {
	sampleLog, _ := ioutil.TempFile("", "sampleLog")
	mw := logger.NewLogMultiWriter(sampleLog.Name())

	server := httptest.NewServer(api.HandleLogSocket(mw))
	defer server.Close()

	dialer := websocket.Dialer{}
	header := make(http.Header)

	wsUrl := strings.ReplaceAll(server.URL+"/api/log/socket", "http", "ws")

	ws, _, err := dialer.Dial(wsUrl, header)
	defer ws.Close()
	assert.NoError(t, err)

	mw.Write([]byte("Hello world!\n"))

	mt, bytes, err := ws.ReadMessage()
	assert.NoError(t, err)

	assert.Equal(t, websocket.TextMessage, mt)
	assert.Equal(t, "Hello world!\n", string(bytes))
}
