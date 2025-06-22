package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Email string

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	
	if email == "" {
		return "", fmt.Errorf("email address is required")
	}
	
	if len(email) > 254 {
		return "", fmt.Errorf("email address is too long")
	}
	
	if !emailRegex.MatchString(email) {
		return "", fmt.Errorf("invalid email format")
	}
	
	return Email(strings.ToLower(email)), nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) IsValid() bool {
	return emailRegex.MatchString(string(e))
}