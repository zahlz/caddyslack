package caddyslack

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteJSONFromReader_NothingToDelete(t *testing.T) {
	bytesIn := []byte(`{"text":"hey"}`)
	readerIn := bytes.NewBuffer(bytesIn)
	pointsToDelete := []string{}

	readerOut, err := deleteJSONFromReader(readerIn, pointsToDelete)

	assert.NoError(t, err)
	bytesOut, err := ioutil.ReadAll(readerOut)
	assert.NoError(t, err)
	assert.Equal(t, bytesOut, bytesIn)
}

func TestDeleteJSONFromReader_InvalidJSON(t *testing.T) {
	readerIn := bytes.NewBuffer([]byte(`]`))
	pointsToDelete := []string{"test"}

	_, err := deleteJSONFromReader(readerIn, pointsToDelete)

	assert.Error(t, err)
}

func TestDeleteJSONFromReader_DeleteExisting(t *testing.T) {
	bytesIn := []byte(`{"channel":"xy","text":"hey"}`)
	readerIn := bytes.NewBuffer(bytesIn)
	pointsToDelete := []string{"channel"}

	readerOut, err := deleteJSONFromReader(readerIn, pointsToDelete)

	assert.NoError(t, err)
	assert.NotEqual(t, readerIn, readerOut)
	bytesOut, err := ioutil.ReadAll(readerOut)
	assert.NoError(t, err)
	bytesExpected := []byte(`{"text":"hey"}`)
	assert.Equal(t, bytesOut, bytesExpected)
}

func TestDeleteJSONFromReader_DeleteNonexisting(t *testing.T) {
	bytesIn := []byte(`{"text":"hey"}`)
	readerIn := bytes.NewBuffer(bytesIn)
	pointsToDelete := []string{"notExsisting"}

	readerOut, err := deleteJSONFromReader(readerIn, pointsToDelete)

	assert.NoError(t, err)
	bytesOut, err := ioutil.ReadAll(readerOut)
	assert.NoError(t, err)
	assert.Equal(t, bytesOut, bytesIn)
}

func TestDeleteJSONFromReader_DeleteNested(t *testing.T) {
	bytesIn := []byte(`{"channel":{"x":"top"},"text":"hey"}`)
	readerIn := bytes.NewBuffer(bytesIn)
	pointsToDelete := []string{"channel.x"}

	readerOut, err := deleteJSONFromReader(readerIn, pointsToDelete)

	assert.NoError(t, err)
	assert.NotEqual(t, readerIn, readerOut)
	bytesOut, err := ioutil.ReadAll(readerOut)
	assert.NoError(t, err)
	bytesExpected := []byte(`{"channel":{},"text":"hey"}`)
	assert.Equal(t, string(bytesOut), string(bytesExpected))
}
