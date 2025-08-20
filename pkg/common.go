package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

func CreateContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return context.WithTimeout(ctx, timeout)
}

func ValidateInput(payload interface{}, validate *validator.Validate) []string {
	if payload == nil {
		return []string{"Invalid Payload"}
	}

	// save error messages here
	var errMessage []string

	errors := validate.Struct(payload)
	if errors != nil {
		// loop through all possible errors,
		// then give appropriate message based on
		// defined error tag, StructField, etc
		for _, err := range errors.(validator.ValidationErrors) {
			if err.Tag() == "required" {
				errMessage = append(errMessage, err.StructField()+" field is required")
				continue
			}

			if err.Tag() == "username" {
				errMessage = append(errMessage, err.StructField()+" contains illegal characters")
				continue
			}

			if err.Tag() == "min" {
				errMessage = append(errMessage, err.StructField()+" field does not meet minimum characters")
				continue
			}

			if err.Tag() == "max" {
				errMessage = append(errMessage, err.StructField()+" field exceed max characters")
				continue
			}

			if err.Tag() == "url" {
				errMessage = append(errMessage, err.StructField()+" field is not a valid URL")
				continue
			}

			if err.Tag() == "email" {
				errMessage = append(errMessage, err.StructField()+" field is not a valid email")
				continue
			}

			if err.Tag() == "oneof" && err.StructField() == "Role" {
				errMessage = append(errMessage, err.StructField()+" field is not a valid role")
				continue
			}

			// raw error which is not covered above
			errMessage = append(errMessage, "Error on field "+err.StructField())
		}
	}

	return errMessage
}

// parseMetadata attempts to decode JSONB data that may be stored either as
// raw JSON (e.g. {"color":"#fff"}) or as a JSON **string** (e.g.
// "{\"color\":\"#fff\"}"). The latter case happens if the value was
// inserted into PostgreSQL wrapped in quotes. We first try to unmarshal
// directly; if that fails we attempt to unmarshal into a string and then parse
// the inner JSON.
func ParseMetadata(raw []byte) (map[string]interface{}, error) {
	var metadata map[string]interface{}
	if raw == nil {
		return metadata, nil
	}

	// Happy path: raw is valid JSON object/array
	if err := json.Unmarshal(raw, &metadata); err == nil {
		return metadata, nil
	}

	// Fallback: raw might be a quoted JSON string
	var quoted string
	if err := json.Unmarshal(raw, &quoted); err == nil {
		if err2 := json.Unmarshal([]byte(quoted), &metadata); err2 == nil {
			return metadata, nil
		}
	}
	return nil, errors.New("invalid metadata JSON format")
}
