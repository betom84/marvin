package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
)

type LogMultiWriter struct {
	mutex       sync.Mutex
	writers     []io.Writer
	multiWriter io.Writer
}

func (l *LogMultiWriter) Append(w io.Writer) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.writers = append(l.writers, w)
	l.multiWriter = io.MultiWriter(l.writers...)
}

func (l *LogMultiWriter) Remove(r io.Writer) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	newWriters := []io.Writer{}

	for _, w := range l.writers {
		if w == r {
			continue
		}

		newWriters = append(newWriters, w)
	}

	l.writers = newWriters
	l.multiWriter = io.MultiWriter(l.writers...)
}

func (l LogMultiWriter) Write(p []byte) (n int, err error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return l.multiWriter.Write(p)
}

func NewLogMultiWriter(logOutput string) *LogMultiWriter {
	var err error
	var base = os.Stdout

	if logOutput != "" && logOutput != "stdout" {
		base, err = os.OpenFile(logOutput, syscall.O_RDWR|syscall.O_CREAT|syscall.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("could not create log '%s'; %s", logOutput, err))
		}
	}

	return &LogMultiWriter{
		writers:     []io.Writer{base},
		multiWriter: io.MultiWriter(base),
	}
}
