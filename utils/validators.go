package utils

import (
	"net/mail"
	"strings"
)

func IsEmpty(field string) bool {
	return len(strings.TrimSpace(field)) == 0
}

func IsNotEmpty(field string) bool {
	return len(strings.TrimSpace(field)) > 0
}

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
