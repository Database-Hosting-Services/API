package utils

import "github.com/google/uuid"

func GenerateOID() string {
	return uuid.New().String() // 36 character
}
