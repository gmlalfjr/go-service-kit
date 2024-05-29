package transform

import (
	"encoding/json"
	"fmt"
)

func InterfaceToBytes(v interface{}) ([]byte, error) {
	var vals = make([]byte, 0)
	switch val := v.(type) {
	case string:
		vals = []byte(val)
	case []byte:
		vals = val
	case int64, int32, int, float64, float32:
		vals = []byte(fmt.Sprintf("%v", val))
	default:
		var err error
		vals, err = json.Marshal(val)
		if err != nil {
			return nil, err
		}
	}

	return vals, nil
}

func InterfaceToString(v interface{}) (string, error) {
	vals, err := InterfaceToBytes(v)
	return string(vals), err
}
