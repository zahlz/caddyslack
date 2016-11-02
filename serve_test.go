package caddyslack

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/stretchr/testify/assert"
)

func newTestHandler(t *testing.T, caddyFile string) *handler {
	c := caddy.NewTestController("http", caddyFile)
	slackConfig, err := parse(c)
	assert.NoError(t, err)
	h := newHandler(slackConfig)
	h.Next = httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
		return http.StatusTeapot, nil
	})
	return h
}

//newTestTarget listens until timeout and returns on the first incoming request
func newTestTarget(addr string, timeoutSec time.Duration) (*http.Request, error) {
	reqCh := make(chan *http.Request)
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		reqCh <- req
	})

	go func() { panic(http.ListenAndServe(addr, nil)) }()

	select {
	case req := <-reqCh:
		return req, nil
	case <-time.After(time.Second * timeoutSec):
		return nil, fmt.Errorf("Timeout after %d seconds", timeoutSec)
	}
}

func TestServeHTTP_ShouldForwardEmptyRequest(t *testing.T) {
	h := newTestHandler(t, `slack {
    url http://localhost:9997
    }`)

	req, err := http.NewRequest("POST", "/slack", nil)
	assert.NoError(t, err)

	go func() {
		requestToSlack, targetErr := newTestTarget("localhost:9997", 1)
		assert.NoError(t, targetErr)
		bodyToSlack, targetErr := ioutil.ReadAll(requestToSlack.Body)
		assert.NoError(t, targetErr)
		assert.Empty(t, bodyToSlack)
	}()

	w := httptest.NewRecorder()
	statusForCaddy, err := h.ServeHTTP(w, req)
	assert.NoError(t, err)
	assert.Exactly(t, statusForCaddy, StatusEmpty)

	bytes, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Exactly(t, w.Code, http.StatusOK, string(bytes))
}

/*func TestServeHTTP_ShouldForwardRequestUnmodified(t *testing.T) {
	h := newTestHandler(t, `slack {
    url http://localhost:9998
    }`)

	jsonStr := []byte(`{"text":"hello"}`)
	req, err := http.NewRequest("POST", "/slack", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)

	go func() {
		requestToSlack, targetErr := newTestTarget("localhost:9998", 1)
		assert.NoError(t, targetErr)
		bodyToSlack, targetErr := ioutil.ReadAll(requestToSlack.Body)
		assert.NoError(t, targetErr)
		assert.Equal(t, bodyToSlack, jsonStr)
	}()

	w := httptest.NewRecorder()
	statusForCaddy, err := h.ServeHTTP(w, req)
	assert.NoError(t, err)
	assert.Exactly(t, statusForCaddy, StatusEmpty)

	bytes, err := ioutil.ReadAll(w.Result().Body)
	assert.NoError(t, err)
	assert.Exactly(t, w.Code, http.StatusOK, string(bytes))
}*/
