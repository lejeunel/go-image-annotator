package event

import (
	"encoding/json"
)

func serialize(m map[string]string) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func deserialize(s string) (map[string]string, error) {
	var m map[string]string

	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
