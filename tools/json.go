package tools

import (
	"encoding/json"
	"strings"
)

// JSONMessage takes a string and returns a JSON formatted string
func JSONMessage(msg string) string {
	return JSONMessageMap(msg, nil)
}

// JSONMessageMap takes a string and a map[string]interface{} and returns a JSON formatted string
func JSONMessageMap(msg string, data map[string]interface{}) string {
	blob := map[string]interface{}{
		"msg":  msg,
		"data": data,
	}

	j, err := json.Marshal(blob)
	if err != nil {
		return `{"msg": "` + strings.Replace(msg, "\"", "\\\"", -1) + `"}`
	}

	return string(j)
}
