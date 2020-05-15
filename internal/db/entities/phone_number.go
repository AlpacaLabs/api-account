package entities

import (
	"time"

	"github.com/rs/xid"

	clock "github.com/AlpacaLabs/go-timestamp"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	"github.com/guregu/null"
)

// PhoneNumber is a representation of a user's email address.
type PhoneNumber struct {
	ID             string
	CreatedAt      time.Time
	LastModifiedAt time.Time
	DeletedAt      null.Time
	Confirmed      bool
	PhoneNumber    string
	AccountID      string
}

type NewPhoneNumberInput struct {
	PhoneNumber string
	AccountID   string
}

func NewPhoneNumber(in NewPhoneNumberInput) PhoneNumber {
	now := time.Now()
	return PhoneNumber{
		ID:             xid.New().String(),
		CreatedAt:      now,
		LastModifiedAt: now,
		DeletedAt:      null.TimeFromPtr(nil),
		PhoneNumber:    in.PhoneNumber,
		AccountID:      in.AccountID,
	}
}

func (e PhoneNumber) ToProtobuf() *accountV1.PhoneNumber {
	return &accountV1.PhoneNumber{
		Id:          e.ID,
		CreatedAt:   clock.TimeToTimestamp(e.CreatedAt),
		Confirmed:   e.Confirmed,
		PhoneNumber: e.PhoneNumber,
		AccountId:   e.AccountID,
	}
}
