package lib

import "github.com/google/uuid"

// constants ...
const (
	ClientPrefix = "$"
)

// NewID ...
func NewID(prefix string) string {
	return ClientPrefix + uuid.New().String()
}
