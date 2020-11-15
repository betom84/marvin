package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"sync"
)

// Configuration holds app configuration properties
type Configuration struct {
	AlexaServerPort    int      `json:"alexaServerPort"`
	UIServerPort       int      `json:"uiServerPort"`
	UIRoot             string   `json:"uiRoot"`
	Log                string   `json:"log"`
	ValidationEnabled  bool     `json:"validationEnabled"`
	ValidationSchema   string   `json:"validationSchema"`
	SSLCertificate     string   `json:"sslCertificate"`
	SSLKey             string   `json:"sslKey"`
	AmazonClientID     string   `json:"amazonClientID"`
	AmazonClientSecret string   `json:"amazonClientSecret"`
	BasicAuthUser      string   `json:"basicAuthUser"`
	BasicAuthPassword  string   `json:"basicAuthPassword"`
	RestrictedUsers    []string `json:"restrictedUsers"`
	Endpoints          string   `json:"endpoints"`
	HomematicHost      string   `json:"homematicHost"`
	PhilipsHueHost     string   `json:"philipsHueHost"`
	PhilipsHueUser     string   `json:"philipsHueUser"`
}

var (
	once       sync.Once
	current    Configuration
	configFile = flag.String("config", "config/config.json", "Path to json encoded config file")
)

func newConfiguration() Configuration {
	return Configuration{
		AlexaServerPort:    6443,
		UIServerPort:       8081,
		ValidationEnabled:  false,
		ValidationSchema:   "./resources/schema.json",
		SSLCertificate:     "./resources/ssl/certificate.pem",
		SSLKey:             "./resources/ssl/private-key.pem",
		Endpoints:          "./config/endpoints.json",
		UIRoot:             "webapp",
		HomematicHost:      "homematic-ccu3",
		PhilipsHueHost:     "philips-hue",
		PhilipsHueUser:     os.Getenv("MARVIN_PHILIPSHUE_USER"),
		AmazonClientID:     os.Getenv("MARVIN_AMAZON_CLIENT_ID"),
		AmazonClientSecret: os.Getenv("MARVIN_AMAZON_CLIENT_SECRET"),
		BasicAuthUser:      os.Getenv("MARVIN_BASIC_AUTH_USER"),
		BasicAuthPassword:  os.Getenv("MARVIN_BASIC_AUTH_PASSWORD"),
		RestrictedUsers:    strings.Split(os.Getenv("MARVIN_RESTRICTED_USERS"), ","),
	}
}

func loadConfiguration(config string) Configuration {
	conf := newConfiguration()

	_, err := os.Stat(config)
	if os.IsNotExist(err) {
		return conf
	}

	cf, err := os.Open(config)
	if err != nil {
		log.Printf("invalid config; %s", err)
		return conf
	}
	defer cf.Close()

	err = json.NewDecoder(cf).Decode(&conf)
	if err != nil {
		log.Printf("could not decode config; %s", err)
		return conf
	}

	return conf
}

// Get current configuration
func Get() Configuration {
	once.Do(func() {
		flag.Parse()
		current = loadConfiguration(*configFile)
	})

	return current
}

// Set current configuration
func Set(c Configuration) {
	once.Do(func() {})
	current = c
}
