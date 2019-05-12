package routes

import (
	"errors"
	"fmt"
)

var (
	errInvalidBody = errors.New("Could not read JSON body")
)

type missingRequiredFieldErr struct {
	field string
}

func (m missingRequiredFieldErr) Error() string {
	return fmt.Sprintf("Missing required field %v", m.field)
}
