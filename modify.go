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

func onlyJSONFromReader(readerIn io.Reader, pointsOnly []string) (readerOut io.Reader, err error) {
	if pointsOnly != nil {
		jsonParsed, err := gabs.ParseJSONBuffer(readerIn)
		if err != nil {
			return readerOut, err
		}
		outContainer := gabs.New()
		for _, point := range pointsOnly {
			if jsonParsed.ExistsP(point) {
				//fmt.Printf("%v: %v\n", point, jsonParsed.Path(point).Data())
				_, err = outContainer.SetP(jsonParsed.Path(point).Data(), point)
				if err != nil {
					return readerOut, err
				}
			}
		}
		//fmt.Printf("container %#v\n", outContainer)
		return bytes.NewReader(outContainer.Bytes()), nil
	}
	return readerIn, nil
}
