package transform

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/skolldire/cash-manager-toolkit/pkg/kit/error_wrapper"
	"net/http"
)

// BytesToModel converts a byte array to a model
func BytesToModel[O any](c []byte) (O, error) {
	h := *new(O)
	e := map[string]interface{}{}
	err := json.Unmarshal(c, &e)
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &h,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	err = decoder.Decode(e)
	if err != nil {
		return h, error_wrapper.NewCommonApiError("TRF-001",
			"[Convert Data To Response]Failed to convert byte array to struct", err, http.StatusInternalServerError)
	}
	return h, nil
}

// ModelToBytes converts a model to a byte array
func ModelToBytes[O any](c O) ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, error_wrapper.NewCommonApiError("TRF-002",
			"[Convert Data To Response]Failed to convert struct to bytes", err, http.StatusInternalServerError)
	}
	return b, nil
}

// MapToStructure converts a map to a structure
func MapToStructure[O any](c map[string]interface{}) (O, error) {
	h := *new(O)
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &h,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	err := decoder.Decode(c)
	if err != nil {
		return h, error_wrapper.NewCommonApiError("TRF-003",
			"[Convert Data To Response]Failed to convert map to struct", err, http.StatusInternalServerError)
	}
	return h, nil
}

// StructToMap converts a struct to a map[string]interface{} respecting the json annotations
// Returns an error if the input is not a struct.
func StructToMap(data interface{}) (map[string]interface{}, error) {

	var mapa map[string]interface{}

	// check if the input value is nil
	if data == nil {
		return mapa, fmt.Errorf("[StructToMap] input value is nil")
	}

	// convert the struct to JSON
	dadosEmJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// deserialize the JSON to the map
	err = json.Unmarshal(dadosEmJson, &mapa)
	if err != nil {
		return nil, err
	}

	return mapa, nil
}
