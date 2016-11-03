package caddyslack

import (
	"bytes"
	"io"

	"github.com/Jeffail/gabs"
)

func deleteJSONFromReader(readerIn io.Reader, pointsToDelete []string) (readerOut io.Reader, err error) {
	if len(pointsToDelete) <= 0 {
		return readerIn, nil
	}

	jsonParsed, err := gabs.ParseJSONBuffer(readerIn)
	if err != nil {
		return readerOut, err
	}
	for _, pointToDelete := range pointsToDelete {
		jsonParsed.DeleteP(pointToDelete)
	}
	return bytes.NewReader(jsonParsed.Bytes()), nil
}

func onlyJSONFromReader(readerIn io.Reader, pointsOnly []string) (readerOut io.Reader, err error) {
	if pointsOnly == nil {
		return readerIn, nil
	}

	jsonParsed, err := gabs.ParseJSONBuffer(readerIn)
	if err != nil {
		return readerOut, err
	}

	outContainer := gabs.New()
	for _, point := range pointsOnly {
		if jsonParsed.ExistsP(point) {
			_, err = outContainer.SetP(jsonParsed.Path(point).Data(), point)
			if err != nil {
				return readerOut, err
			}
		}
	}
	return bytes.NewReader(outContainer.Bytes()), nil
}
