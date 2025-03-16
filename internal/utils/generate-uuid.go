package utils

import "github.com/google/uuid"

func GenerateRandomId() string {
	return uuid.NewString()
}
