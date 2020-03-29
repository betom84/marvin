package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"marvin/config"
)

func HandleLogGet() func(w http.ResponseWriter, r *http.Request) error {

	var response = struct {
		Lines []string `json:"lines"`
	}{}

	return func(w http.ResponseWriter, r *http.Request) error {

		offset, _ := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 0)
		limit, _ := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 0)

		if limit <= 0 {
			limit = 10
		}

		lines, err := getLatestLogMessages(uint(limit), uint(offset))
		if err != nil {
			return err
		}

		response.Lines = lines
		json.NewEncoder(w).Encode(response)
		w.Header().Set("Content-Type", "application/json")

		return nil
	}

}

func getLatestLogMessages(limit, offset uint) ([]string, error) {

	logFilePath := config.Get().Log
	if logFilePath == "stdout" {
		return []string{}, nil
	}

	logfile, err := os.Open(logFilePath)
	if err != nil {
		return nil, err
	}
	defer logfile.Close()

	var lineCount int
	var neededCount = int(limit + offset)
	var seekLength int64 = 1024

	currentOffset, _ := logfile.Seek(0, io.SeekEnd)

	for lineCount <= neededCount && currentOffset > 0 {
		if seekLength > currentOffset {
			seekLength = currentOffset
		}

		currentOffset, _ = logfile.Seek(seekLength*-1, io.SeekCurrent)

		readBuffer := make([]byte, seekLength)
		readCount, err := logfile.Read(readBuffer)
		if err != nil {
			return nil, err
		}

		lineCount += bytes.Count(readBuffer, []byte("\n"))
		currentOffset, _ = logfile.Seek(int64(readCount)*-1, io.SeekCurrent)
	}

	if uint(lineCount) <= offset {
		return []string{}, nil
	}

	entries := []string{}
	scanner := bufio.NewScanner(logfile)
	for scanner.Scan() == true {
		entries = append(entries, scanner.Text())
	}

	sliceStart := len(entries) - neededCount
	if sliceStart < 0 {
		sliceStart = 0
	}

	return entries[sliceStart : len(entries)-int(offset)], nil
}
