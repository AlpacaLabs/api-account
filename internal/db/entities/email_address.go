package entities

import (
	clock "github.com/AlpacaLabs/go-timestamp"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	"github.com/guregu/null"
)

// EmailAddress is a representation of a user's email address.
type EmailAddress struct {
	ID           string
	Created      null.Time
	Deleted      null.Time
	LastModified null.Time
	Confirmed    bool
	Primary      bool
	EmailAddress string
	AccountID    string
}

func (e EmailAddress) ToProtobuf() *accountV1.EmailAddress {
	return &accountV1.EmailAddress{
		Id:             e.ID,
		CreatedAt:      clock.TimeToTimestamp(e.Created.ValueOrZero()),
		LastModifiedAt: clock.TimeToTimestamp(e.LastModified.ValueOrZero()),
		Deleted:        e.Deleted.Valid,
		DeletedAt:      clock.TimeToTimestamp(e.Deleted.ValueOrZero()),
		Confirmed:      e.Confirmed,
		Primary:        e.Primary,
		EmailAddress:   e.EmailAddress,
		AccountId:      e.AccountID,
	}
}
