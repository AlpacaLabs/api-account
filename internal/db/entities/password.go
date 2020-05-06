package entities

import (
	"github.com/guregu/null"
)

// Password is a representation of a user's password.
type Password struct {
	Id             string    `json:"id"`
	Created        null.Time `json:"created_at"`
	IterationCount int
	Salt           []byte
	PasswordHash   []byte
	AccountID      string `json:"account_id"`
}
