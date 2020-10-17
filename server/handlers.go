package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gorilla/mux"
)

const (
	B float64 = 1 << (10 * iota)
	KB
	MB
	GB
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"status\": \"ok\"}"))
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	statusCode, err := strconv.Atoi(vars["status"])
	if err != nil || statusCode < 400 || statusCode >= 600 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func newDataHandler(maxDataSize string) (*dataHandler, error) {
	dh := &dataHandler{}
	mds, err := dh.parseBytes(maxDataSize)
	dh.maxDataSize = mds
	return dh, err
}

type dataHandler struct {
	maxDataSize uint64
}

func (dh *dataHandler) parseBytes(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	split := strings.IndexFunc(s, unicode.IsLetter)
	if split == -1 {
		return 0, errors.New("invalid data size format")
	}

	byteString, unit := s[:split], s[split:]
	bytes, err := strconv.ParseFloat(byteString, 64)
	if err != nil {
		return 0, errors.New("invalid data size")
	}

	switch unit {
	case "B":
		return uint64(bytes * B), nil
	case "KB":
		return uint64(bytes * KB), nil
	case "MB":
		return uint64(bytes * MB), nil
	case "GB":
		return uint64(bytes * GB), nil
	default:
		return 0, errors.New("invalid data unit")
	}
}

func (dh *dataHandler) Handler(w http.ResponseWriter, r *http.Request) {
	size, err := dh.parseBytes(r.Form.Get("s"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if size > dh.maxDataSize {
		http.Error(w, "size exceeds maximum", http.StatusBadRequest)
		return
	}
	data := make([]byte, size)
	w.Write(data)
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
