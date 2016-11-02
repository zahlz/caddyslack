package caddyslack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/zahlz/caddyslack/bufpool"
)

// StatusEmpty returned by mailout middleware because the proper status gets
// written previously
const StatusEmpty = 0

const (
	headerContentType         = "Content-Type"
	headerApplicationJSONUTF8 = "application/json; charset=utf-8"
)

type handler struct {
	Next   httpserver.Handler
	config *config
}

func newHandler(sc *config) *handler {
	fmt.Printf("%+v\n", sc)
	return &handler{
		config: sc,
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// Validate
	if r.URL.Path != h.config.endpoint {
		return h.Next.ServeHTTP(w, r)
	}
	if r.Method != "POST" {
		return h.writeJSON(JSONError{
			Code:  http.StatusMethodNotAllowed,
			Error: http.StatusText(http.StatusMethodNotAllowed),
		}, w)
	}

	// Modify
	reader := printReader(r.Body)
	reader, err := deleteJSONFromReader(reader, h.config.delete)
	if err != nil {
		return h.writeJSON(JSONError{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		}, w)
	}
	reader = printReader(reader)
	reader, err = onlyJSONFromReader(reader, h.config.only)
	if err != nil {
		return h.writeJSON(JSONError{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		}, w)
	}
	reader = printReader(reader)
	// Proxy
	res, err := http.Post(h.config.remoteURL, "application/json", reader)
	if err != nil {
		return h.writeJSON(JSONError{
			Code:  http.StatusTeapot,
			Error: err.Error(),
		}, w)
	}

	w.Header().Set(headerContentType, headerApplicationJSONUTF8)
	w.WriteHeader(res.StatusCode)

	return 0, err
}

// JSONError defines how an REST JSON looks like.
// Code 200 and empty Error specifies a successful request
// Any other Code value s an error.
type JSONError struct {
	// Code represents the HTTP Status Code, a work around.
	Code int `json:"code,omitempty"`
	// Error the underlying error, if there is one.
	Error string `json:"error,omitempty"`
}

func (h *handler) writeJSON(je JSONError, w http.ResponseWriter) (int, error) {
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	w.Header().Set(headerContentType, headerApplicationJSONUTF8)

	// https://github.com/mholt/caddy/issues/637#issuecomment-189599332
	w.WriteHeader(je.Code)

	if err := json.NewEncoder(buf).Encode(je); err != nil {
		return http.StatusInternalServerError, err
	}
	if _, err := w.Write(buf.Bytes()); err != nil {
		return http.StatusInternalServerError, err
	}

	return StatusEmpty, nil
}

func printReader(reader io.Reader) io.Reader {
	if reader != nil {
		allBytes, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Printf("Error reading reader: %v\n", err)
		}
		fmt.Println(string(allBytes))
		return bytes.NewBuffer(allBytes)
	} else {
		fmt.Println("reader == nil")
	}
	return reader
}
