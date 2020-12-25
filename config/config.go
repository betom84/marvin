package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
)

// Configuration holds app configuration properties
type Configuration struct {
	AlexaServerPort    int    `json:"alexaServerPort"`
	UIServerPort       int    `json:"uiServerPort"`
	UIRoot             string `json:"uiRoot"`
	Log                string `json:"log"`
	ValidationEnabled  bool   `json:"validationEnabled"`
	ValidationSchema   string `json:"validationSchema"`
	SSLCertificate     string `json:"sslCertificate"`
	SSLKey             string `json:"sslKey"`
	AmazonClientID     string `json:"amazonClientID"`
	AmazonClientSecret string `json:"amazonClientSecret"`
	BasicAuthUser      string `json:"basicAuthUser"`
	BasicAuthPassword  string `json:"basicAuthPassword"`
	RestrictedUser     string `json:"restrictedUser"`
	Endpoints          string `json:"endpoints"`
	HomematicHost      string `json:"homematicHost"`
	PhilipsHueHost     string `json:"philipsHueHost"`
	PhilipsHueUser     string `json:"philipsHueUser"`
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
		PhilipsHueUser:     "$MARVIN_PHILIPSHUE_USER",
		AmazonClientID:     "$MARVIN_AMAZON_CLIENT_ID",
		AmazonClientSecret: "$MARVIN_AMAZON_CLIENT_SECRET",
		BasicAuthUser:      "$MARVIN_BASIC_AUTH_USER",
		BasicAuthPassword:  "$MARVIN_BASIC_AUTH_PASSWORD",
		RestrictedUser:     "$MARVIN_RESTRICTED_USER",
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

func resolveEnvironmentVariables(conf Configuration) Configuration {
	// https://blog.golang.org/laws-of-reflection
	confValues := reflect.ValueOf(&conf).Elem()

	for i := 0; i < confValues.NumField(); i++ {
		fieldValue := confValues.Field(i).Interface()

		if _, ok := fieldValue.(string); ok != true {
			continue
		}

		stringValue := fieldValue.(string)
		if strings.HasPrefix(stringValue, "$") != true {
			continue
		}

		confValues.Field(i).SetString(os.Getenv(stringValue[1:]))
	}

	return conf
}

// Get current configuration
func Get() Configuration {
	return resolveEnvironmentVariables(GetRaw())
}

// GetRaw like Get but without resolving environment variables
func GetRaw() Configuration {
	once.Do(func() {
		flag.Parse()
		current = loadConfiguration(*configFile)
	})

	return current
}

// GetDefault configuration without resolving environment variables
func GetDefault() Configuration {
	return newConfiguration()
}

// Set current configuration
func Set(c Configuration) {
	once.Do(func() {})
	current = c
}
