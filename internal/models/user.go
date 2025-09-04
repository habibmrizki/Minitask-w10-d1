package models

import (
	"errors"
	"regexp"
)

var (
	emailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	lowercaseRegex   = regexp.MustCompile(`[a-z]`)
	uppercaseRegex   = regexp.MustCompile(`[A-Z]`)
	numberRegex      = regexp.MustCompile(`[0-9]`)
	specialCharRegex = regexp.MustCompile(`[!@#$%^&*()_+=\-<>?|~{}]`)
)

var users = make(map[string]User)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error"`
	UserID  int    `json:"user_id,omitempty"`
}

func InitUsers() {

}

func GetUsersMap() map[string]User {
	return users
}

// ValidateEmail memeriksa format email
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("format email tidak valid")
	}
	return nil
}

// ValidateUserCredentials memeriksa validitas password
func ValidateUserCredentials(user User) error {
	if !lowercaseRegex.MatchString(user.Password) {
		return errors.New("kata sandi harus mengandung setidaknya satu huruf kecil")
	}
	if !uppercaseRegex.MatchString(user.Password) {
		return errors.New("kata sandi harus mengandung setidaknya satu huruf besar")
	}
	if !numberRegex.MatchString(user.Password) {
		return errors.New("kata sandi harus mengandung setidaknya satu angka")
	}
	if !specialCharRegex.MatchString(user.Password) {
		return errors.New("kata sandi harus mengandung setidaknya satu karakter khusus")
	}
	return nil
}
