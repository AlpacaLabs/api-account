package entities

import (
	"time"

	clock "github.com/AlpacaLabs/go-timestamp"
	clocksql "github.com/AlpacaLabs/go-timestamp-sql"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	"github.com/guregu/null"
)

// EmailAddress is a representation of a user's email address.
type EmailAddress struct {
	ID             string
	CreatedAt      time.Time
	LastModifiedAt time.Time
	DeletedAt      null.Time
	Confirmed      bool
	Primary        bool
	EmailAddress   string
	AccountID      string
}

type NewEmailAddressInput struct {
	Primary      bool
	EmailAddress string
	AccountID    string
}

func NewEmailAddress(in NewEmailAddressInput) EmailAddress {
	now := time.Now()
	return EmailAddress{
		ID:             "",
		CreatedAt:      now,
		LastModifiedAt: now,
		DeletedAt:      null.TimeFromPtr(nil),
		Primary:        in.Primary,
		EmailAddress:   in.EmailAddress,
		AccountID:      in.AccountID,
	}
}

func (e EmailAddress) ToProtobuf() *accountV1.EmailAddress {
	return &accountV1.EmailAddress{
		Id:             e.ID,
		CreatedAt:      clock.TimeToTimestamp(e.CreatedAt),
		LastModifiedAt: clock.TimeToTimestamp(e.LastModifiedAt),
		Deleted:        e.DeletedAt.Valid,
		DeletedAt:      clocksql.TimestampFromNullTime(e.DeletedAt),
		Confirmed:      e.Confirmed,
		Primary:        e.Primary,
		EmailAddress:   e.EmailAddress,
		AccountId:      e.AccountID,
	}
}
