package entities

import (
	"time"

	"github.com/guregu/null"
)

type Account struct {
	ID                    string
	CreatedAt             time.Time
	LastModifiedAt        time.Time
	DeletedAt             null.Time
	Username              null.String
	CurrentPasswordID     null.String
	PrimaryEmailAddressID null.String
}
