package util

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

func GenerateClientId() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}
