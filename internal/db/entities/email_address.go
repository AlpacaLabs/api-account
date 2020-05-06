package entities

import (
	clock "github.com/AlpacaLabs/go-timestamp"
	authV1 "github.com/AlpacaLabs/protorepo-auth-go/alpacalabs/auth/v1"
	"github.com/guregu/null"
)

// TODO add newEmailAddress method
// TODO manually manage timestamps
// TODO return error if PUT/POST includes more fields than it should? or do we quietly process only the fields we care about

// EmailAddress is a representation of a user's email address.
type EmailAddress struct {
	ID           string    `json:"id"`
	Created      null.Time `json:"created_at"`
	Deleted      null.Time `json:"deleted_at"`
	LastModified null.Time `json:"last_modified_at"`
	Confirmed    bool      `json:"confirmed"`
	Primary      bool      `json:"primary"`
	EmailAddress string    `json:"email_address"`
	AccountID    string    `json:"account_id"`
}

func (e EmailAddress) ToProtobuf() *authV1.EmailAddress {
	return &authV1.EmailAddress{
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
