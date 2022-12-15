package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"marvin/alexa"
	"marvin/config"
	"marvin/devices/homematic"
	"marvin/devices/hue"
	"marvin/devices/shelly"
	"marvin/logger"
	"marvin/ui"

	"github.com/betom84/go-alexa/smarthome"
)

func main() {
	writer := logger.NewLogMultiWriter(config.Get().Log)

	log.SetFlags(log.LstdFlags)
	log.SetOutput(writer)
	smarthome.Log = smarthome.NewDefaultLogger(smarthome.Trace, writer)

	var sigChan = make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT)

	alexaServer := newAlexaServer()
	err := alexaServer.Start()
	if err != nil {
		log.Printf("failed to start alexa server; %s\n", err)
	}

	startUIServer(alexaServer, writer)

	sig := <-sigChan
	log.Printf("exiting due of signal %+v...\n", sig)
	os.Exit(0)
}

func newAlexaServer() *alexa.Server {
	config := config.Get()

	server := alexa.NewServer()
	server.Addr = fmt.Sprintf(":%d", config.AlexaServerPort)
	server.CertFile = config.SSLCertificate
	server.KeyFile = config.SSLKey
	server.AmazonClientID = config.AmazonClientID
	server.AmazonClientSecret = config.AmazonClientSecret
	server.BasicAuthUser = config.BasicAuthUser
	server.BasicAuthPassword = config.BasicAuthPassword
	server.ValidationEnabled = config.ValidationEnabled
	server.ValidationSchema = config.ValidationSchema
	server.RestrictedUsers = config.RestrictedUsers
	server.Endpoints = config.Endpoints

	server.NewDeviceFunc("homematic", func(id string) (interface{}, error) {
		return homematic.NewDevice(id, config.HomematicHost)
	})

	server.NewDeviceFunc("hue", func(id string) (interface{}, error) {
		return hue.NewLight(id, config.PhilipsHueHost, config.PhilipsHueUser)
	})

	server.NewDeviceFunc("shelly", func(id string) (interface{}, error) {
		return shelly.NewDevice(id)
	})

	return server
}

func startUIServer(a *alexa.Server, mw *logger.LogMultiWriter) {
	server := ui.NewServer(a, mw)

	go func() {
		log.Printf("starting ui on :%d\n", config.Get().UIServerPort)

		err := http.ListenAndServe(fmt.Sprintf(":%d", config.Get().UIServerPort), server)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}
