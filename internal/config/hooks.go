package config

import (
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

func DecodeLogTypeHookFunc() mapstructure.DecodeHookFuncType {
	// Wrapped in a function call to add optional input parameters (eg. separator)
	return func(
		f reflect.Type, // data type
		t reflect.Type, // target data type
		data any, // raw data
	) (any, error) {
		// Check if the data type matches the expected one
		if f.Kind() != reflect.String {
			return data, nil
		}

		// Check if the target type matches the expected one
		if t != reflect.TypeFor[LogType]() {
			return data, nil
		}
		// Format/decode/parse the data and return the new value
		switch data.(string) {
		case "json":
			return LOGTYPE_JSON, nil
		case "pretty":
			return LOGTYPE_PRETTY, nil
		default:
			return nil, fmt.Errorf("unknown log type: %s", data)
		}
	}
}
