package utils

import "encoding/json"

func JsonToMap(jsonString string) ([]map[string]interface{}, error) {
	m := make([]map[string]interface{}, 999)
	err := json.Unmarshal([]byte(jsonString), &m)
	return m, err
}
