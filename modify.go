package caddyslack

import (
	"bytes"
	"io"

	"github.com/Jeffail/gabs"
)

func deleteJSONFromReader(readerIn io.Reader, pointsToDelete []string) (readerOut io.Reader, err error) {
	if len(pointsToDelete) > 0 {
		jsonParsed, err := gabs.ParseJSONBuffer(readerIn)
		if err != nil {
			return readerOut, err
		}
		for _, pointToDelete := range pointsToDelete {
			jsonParsed.DeleteP(pointToDelete)
		}
		return bytes.NewReader(jsonParsed.Bytes()), nil
	}
	return readerIn, nil
}
