package uuid

import (
	"github.com/google/uuid"
)

// Generate - uuid generator
func Generate() string {
	return uuid.New().String()
}
