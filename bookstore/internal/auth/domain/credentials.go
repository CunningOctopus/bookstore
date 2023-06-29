package domain

import (
	"errors"
	"net/mail"
	"unicode"
)

type Credentials struct {
	EMail    string `json:"email"`
	Password string `json:"password"`
}

var ErrBadEMail = errors.New("email address must be valid")
var ErrBadPassword = errors.New("password must be at least 8 symbols long and consist of at least digits and letters")

func (creds *Credentials) Validate() error {

	if _, err := mail.ParseAddress(creds.EMail); err != nil {
		return ErrBadEMail
	}

	if len(creds.Password) < 8 || len(creds.Password) > 32 {
		return ErrBadPassword
	}

	var hasLetters, hasDigits bool
	for _, char := range creds.Password {
		if unicode.IsLetter(char) {
			hasLetters = true
		}

		if unicode.IsDigit(char) {
			hasDigits = true
		}

		if hasDigits && hasLetters {
			return nil
		}
	}

	return ErrBadPassword
}
