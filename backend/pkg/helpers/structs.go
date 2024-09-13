package helpers

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func PopulateStruct(payload map[string]interface{}, target interface{}) error {
	return mapstructure.Decode(payload, target)
}

func PrettyStruct(payload interface{}) string {
	s, _ := json.MarshalIndent(payload, "", "  ")
	return string(s)
}

var validate = validator.New()

func ValidateStruct[T any](payload T) []*ValidationErrors {
	var errors []*ValidationErrors
	if err := validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, &ValidationErrors{
				Field: err.Field(),
				Rule:  err.Tag(),
				Value: err.Param(),
			})
		}
	}
	return errors
}

type ValidationErrors struct {
	Field string `json:"field"`
	Rule  string `json:"tag"`
	Value string `json:"value"`
}