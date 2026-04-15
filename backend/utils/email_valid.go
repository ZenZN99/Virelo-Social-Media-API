package utils

import (
	"net/mail"
	"strings"
)

func IsValidEmail(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))

	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	if strings.Contains(email, "..") {
		return false
	}

	if strings.HasPrefix(email, "!") || strings.HasPrefix(email, ".") {
		return false
	}

	return true
}