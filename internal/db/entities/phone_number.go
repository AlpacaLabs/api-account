package entities

import (
	"time"

	clock "github.com/AlpacaLabs/go-timestamp"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	"github.com/guregu/null"
)

// PhoneNumber is a representation of a user's email address.
type PhoneNumber struct {
	ID           string
	Created      time.Time
	LastModified time.Time
	Deleted      null.Time
	Confirmed    bool
	PhoneNumber  string
	AccountID    string
}

func (e PhoneNumber) ToProtobuf() *accountV1.PhoneNumber {
	return &accountV1.PhoneNumber{
		Id:          e.ID,
		CreatedAt:   clock.TimeToTimestamp(e.Created),
		Confirmed:   e.Confirmed,
		PhoneNumber: e.PhoneNumber,
		AccountId:   e.AccountID,
	}
}