package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	assert.NoError(t, os.Setenv("MARVIN_BASIC_AUTH_USER", "chuck.norris"))
	assert.NoError(t, os.Setenv("MARVIN_BASIC_AUTH_PASSWORD", ""))

	*configFile = writeConfig(t, []byte(`{"amazonClientID": "clientID"}`))
	defer os.Remove(*configFile)

	// consider value from config
	assert.Equal(t, "clientID", Get().AmazonClientID)

	// consider defaults
	assert.Equal(t, 8081, Get().UIServerPort)
	assert.Equal(t, 6443, Get().AlexaServerPort)

	// consider config from env
	assert.Equal(t, "chuck.norris", Get().BasicAuthUser)
	assert.Equal(t, "", Get().BasicAuthPassword)

	// raw does not resolve env
	assert.Equal(t, "$MARVIN_BASIC_AUTH_USER", GetRaw().BasicAuthUser)
	assert.Equal(t, "$MARVIN_BASIC_AUTH_PASSWORD", GetRaw().BasicAuthPassword)

	// default does not resolve env
	assert.Equal(t, "$MARVIN_BASIC_AUTH_USER", GetDefault().BasicAuthUser)
	assert.Equal(t, "$MARVIN_BASIC_AUTH_PASSWORD", GetDefault().BasicAuthPassword)
}

func writeConfig(t *testing.T, config []byte) string {
	tmpfile, err := ioutil.TempFile("", "example")
	assert.NoError(t, err)

	_, err = tmpfile.Write(config)
	assert.NoError(t, err)

	err = tmpfile.Close()
	assert.NoError(t, err)

	return tmpfile.Name()
}
