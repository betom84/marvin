package logger_test

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"marvin/logger"

	"github.com/stretchr/testify/assert"
)

func TestLogMultiWriter(t *testing.T) {
	sampleFile, err := ioutil.TempFile("", "sample")
	assert.NoError(t, err)

	mw := logger.NewLogMultiWriter(sampleFile.Name())

	b := new(bytes.Buffer)
	mw.Append(b)

	_, err = mw.Write([]byte("Hello world!"))
	assert.NoError(t, err)

	assert.Equal(t, "Hello world!", b.String())

	b.Reset()
	mw.Remove(b)

	_, err = mw.Write([]byte("Bye world!"))

	assert.Empty(t, b.Len())
}

func TestLogMultiWriterWithPipe(t *testing.T) {
	sampleFile, err := ioutil.TempFile("", "sample")
	assert.NoError(t, err)

	mw := logger.NewLogMultiWriter(sampleFile.Name())

	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()

	mw.Append(pw)

	r := bufio.NewReader(pr)

	go func() {
		mw.Write([]byte("Hello "))
		mw.Write([]byte("world"))
		mw.Write([]byte("!"))
		mw.Write([]byte("\n"))
	}()

	rb, err := r.ReadBytes('\n')
	assert.NoError(t, err)

	assert.Equal(t, "Hello world!\n", string(rb))
}
