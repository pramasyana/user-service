package helper

import (
	"github.com/google/jsonapi"
)

// MarshalConvertOnePayload function to convert struct response to jsonapi.OnePayLoad so that we can add meta or link data
func MarshalConvertOnePayload(structResponse interface{}) (payload *jsonapi.OnePayload, err error) {
	// set response marshal jsonapi struct
	p, err := jsonapi.Marshal(structResponse)
	if err != nil {

		return nil, err
	}

	var ok bool
	if payload, ok = p.(*jsonapi.OnePayload); !ok {
		return nil, err
	}

	return
}
