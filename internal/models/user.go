package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  Password  `json:"password"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
}

type Password struct {
	Plaintext *string `json:"plaintext"`
	Hash      []byte  `json:"-"`
}

// The Set method calculates the bcrypt Hash of a Plaintext password, and stores both // the Hash and the Plaintext versions in the struct.
func (p *Password) Set(PlaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(PlaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &PlaintextPassword
	p.Hash = hash
	return nil
}

// The Matches method checks whether the provided Plaintext password matches the // Hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *Password) Matches(PlaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(PlaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
