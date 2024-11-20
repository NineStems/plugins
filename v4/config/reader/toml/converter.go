package toml

import (
	"encoding/json"
)

func convertValue(b []byte) any {
	var data any
	if err := json.Unmarshal(b, &data); err != nil {
		return nil
	}

	switch v := data.(type) {
	case float64:
		if v == float64(int(v)) {
			return int(v)
		}
		return v
	case string:
		return v
	case bool:
		return v
	case []interface{}:
		return v
	case map[string]interface{}:
		return v
	default:
		return nil
	}
}
