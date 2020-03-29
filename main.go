package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"marvin/alexa"
	"marvin/config"
	"marvin/devices/homematic"
	"marvin/devices/hue"
	"marvin/ui"

	"github.com/betom84/go-alexa/smarthome"
)

func init() {
	initLog()
}

func initLog() {
	var err error
	var writer io.Writer = os.Stdout

	logOutput := config.Get().Log
	if logOutput != "" && logOutput != "stdout" {
		writer, err = os.OpenFile(logOutput, syscall.O_RDWR|syscall.O_CREAT|syscall.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("could not create log '%s'; %s", logOutput, err))
		}
	}

	log.SetFlags(log.LstdFlags)
	log.SetOutput(writer)

	smarthome.Log = smarthome.NewDefaultLogger(smarthome.Trace, writer)
}

func main() {
	var sigChan = make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT)

	alexaServer := newAlexaServer()
	err := alexaServer.Start()
	if err != nil {
		log.Printf("failed to start alexa server; %s\n", err)
	}

	startUIServer(alexaServer)

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

	server.NewDeviceFunc("homematic", func(id int) interface{} {
		return homematic.NewDevice(id, config.HomematicHost)
	})

	server.NewDeviceFunc("hue", func(id int) interface{} {
		return hue.NewLight(id, config.PhilipsHueHost, config.PhilipsHueUser)
	})

	return server
}

func startUIServer(a *alexa.Server) {
	server := ui.NewServer(a)

	go func() {
		log.Printf("starting ui on :%d\n", config.Get().UIServerPort)

		err := http.ListenAndServe(fmt.Sprintf(":%d", config.Get().UIServerPort), server)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}
