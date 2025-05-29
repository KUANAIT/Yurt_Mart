package utils

import (
	"encoding/json"
	"errors"
)

func StructToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &result)
	return result, err
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	return nil
}
