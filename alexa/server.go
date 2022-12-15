package alexa

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"marvin/metrics"
	"net/http"
	"os"
	"time"

	"github.com/betom84/go-alexa/smarthome"
	"github.com/betom84/go-alexa/smarthome/common/discoverable"
	"github.com/betom84/go-alexa/smarthome/validator"
)

type DeviceHandlerFunc func(id string) (interface{}, error)

// Server to process alexa directives
type Server struct {
	httpServer    *http.Server
	deviceHandler map[string]DeviceHandlerFunc

	Addr               string
	CertFile           string
	KeyFile            string
	ValidationEnabled  bool
	ValidationSchema   string
	AmazonClientID     string
	AmazonClientSecret string
	BasicAuthUser      string
	BasicAuthPassword  string
	RestrictedUsers    []string
	Endpoints          string
}

func NewServer() *Server {
	s := Server{
		deviceHandler: make(map[string]DeviceHandlerFunc),
	}

	return &s
}

// Start alexa server
func (server *Server) Start() error {
	if server.IsRunning() {
		return fmt.Errorf("alexa server already running")
	}

	smarthomeHandler, err := server.newSmarthomeHandler()
	if err != nil {
		return err
	}

	handler := http.NewServeMux()
	handler.Handle("/", metrics.Middleware(smarthomeHandler))

	server.httpServer = &http.Server{
		Addr:    server.Addr,
		Handler: handler,
	}

	go func() {
		log.Printf("starting AlexaServer on %s\n", server.httpServer.Addr)

		err := server.httpServer.ListenAndServeTLS(server.CertFile, server.KeyFile)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return nil
}

func (server Server) newSmarthomeHandler() (*smarthome.Handler, error) {
	if server.AmazonClientID == "" || server.AmazonClientSecret == "" {
		return nil, fmt.Errorf("amazon client credentials missing")
	}

	discoverableEndpoints, err := server.readEndpoints()
	if err != nil {
		return nil, fmt.Errorf("unable to read alexa device endpoints; %s", err)
	}

	authority := smarthome.Authority{
		ClientID:        server.AmazonClientID,
		ClientSecret:    server.AmazonClientSecret,
		RestrictedUsers: server.RestrictedUsers,
	}

	handler := smarthome.NewDefaultHandler(&authority, discoverableEndpoints)
	handler.BasicAuth.Username = server.BasicAuthUser
	handler.BasicAuth.Password = server.BasicAuthPassword
	handler.DeviceFactory = server

	if server.ValidationEnabled && server.ValidationSchema != "" {
		handler.Validator = &validator.Validator{SchemaReference: server.ValidationSchema}
	}

	return handler, nil
}

func (server Server) readEndpoints() ([]discoverable.Endpoint, error) {
	var discoverableEndpoints = []discoverable.Endpoint{}

	_, err := os.Stat(server.Endpoints)
	if os.IsNotExist(err) {
		return discoverableEndpoints, nil
	}

	endpoints, err := os.Open(server.Endpoints)
	if err != nil {
		return nil, err
	}

	ep, err := ioutil.ReadAll(endpoints)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(ep, &discoverableEndpoints)
	if err != nil {
		return nil, err
	}

	return discoverableEndpoints, nil
}

// Stop alexa server
func (server *Server) Stop() error {
	if !server.IsRunning() {
		return fmt.Errorf("alexa server not running")
	}

	server.httpServer.SetKeepAlivesEnabled(false)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := server.httpServer.Shutdown(ctx)
	if err == nil {
		server.httpServer = nil
	}

	log.Printf("stopped AlexaServer\n")

	return err
}

// IsRunning ...
func (server Server) IsRunning() bool {
	return server.httpServer != nil
}

func (server *Server) NewDeviceFunc(epType string, f DeviceHandlerFunc) {
	server.deviceHandler[epType] = f
}

func (server Server) NewDevice(epType string, id string) (interface{}, error) {
	h, ok := server.deviceHandler[epType]
	if !ok {
		return nil, fmt.Errorf("unable to instantiate device for endpoint type %s, id %s", epType, id)
	}

	return h(id)
}
