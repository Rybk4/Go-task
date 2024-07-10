package validation

import (
	"errors"
	"fmt"
	"strconv"
)

func required(fieldName string, value interface{}, arg string) (valid bool, err error) {
	if str, ok := value.(string); ok {
		valid = str != ""
		if !valid {
			err = fmt.Errorf("A value is required")
		}
	} else {
		err = errors.New("The required validator is for strings")
	}
	return
}

func min(fieldName string, value interface{}, arg string) (valid bool, err error) {
	minVal, convErr := strconv.Atoi(arg)
	if convErr != nil {
		panic("Invalid arguments for validator: " + arg)
	}
	err = fmt.Errorf("The minimum value is %v", minVal)

	switch v := value.(type) {
	case int:
		valid = v >= minVal
	case float64:
		valid = v >= float64(minVal)
	case string:
		err = fmt.Errorf("The minimum length is %v characters", minVal)
		valid = len(v) >= minVal
	default:
		err = errors.New("The min validator is for int, float64, and string values")
	}
	return
}
