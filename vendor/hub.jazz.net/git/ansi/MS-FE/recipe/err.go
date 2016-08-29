package recipe

import "fmt"

type Error int

// Error returns the string representation of Error
func (e Error) Error() string {
	if errorMessages[e] == "" {
		return fmt.Sprintf("Error %d unknown", int(e))
	}

	return errorMessages[e]
}

// Error codes. Casted to Error, these yield all possible errors this package
// provides. Use recipe.Error(errorcode).Error() to get a descriptive string for an
// error code.
const (
	E_SENSOR_NOT_FOUND Error = iota
	E_SENSOR_DISABLED
	E_NO_NEXT_OPERATION
	E_RECIPE_TYPE_UNKNOWN
)

var errorMessages = map[Error]string{
	E_SENSOR_NOT_FOUND:    "sensor not found in recipe",
	E_SENSOR_DISABLED:     "sensor disabled",
	E_NO_NEXT_OPERATION:   "no more operations in recipe",
	E_RECIPE_TYPE_UNKNOWN: "recipe type unknown",
}
