package utils

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateTicketNumber membuat UUID dan mengambil 8 karakter pertama
func GenerateTicketNumber(leng int) string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")[:leng]
}
