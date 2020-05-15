package service

import (
	"context"
	"fmt"

	"github.com/AlpacaLabs/api-account/internal/db"
	"github.com/AlpacaLabs/api-account/internal/db/entities"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	"github.com/badoux/checkmail"
	"github.com/rs/xid"
	"github.com/ttacon/libphonenumber"
)

const (
	MinUsernameLength = 4
	MaxUsernameLength = 25
)

var (
	ErrUsernameInvalidLength = fmt.Errorf("username must be between %d and %d characters long", MinUsernameLength, MaxUsernameLength)
)

func (s Service) CreateAccount(ctx context.Context, request *accountV1.CreateAccountRequest) (*accountV1.CreateAccountResponse, error) {
	emailAddress := request.EmailAddress
	username := request.Username
	phoneNumber := request.PhoneNumber

	// Validate email address
	if err := checkmail.ValidateFormat(emailAddress); err != nil {
		return nil, err
	}

	// Validate username
	if username != "" {
		if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
			return nil, ErrUsernameInvalidLength
		}
	}

	// Validate phone number
	if phoneNumber != "" {
		if _, err := libphonenumber.Parse(phoneNumber, "US"); err != nil {
			return nil, err
		}
	}

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		// TODO check if (confirmed or not) email address already exists
		// TODO check if (confirmed or not) phone number already exists
		// TODO check if account already exists with username

		accountID := xid.New().String()

		if err := tx.CreateAccount(ctx, accountID, username); err != nil {
			return err
		}

		if err := tx.CreatePhoneNumber(ctx, entities.NewPhoneNumber(entities.NewPhoneNumberInput{
			PhoneNumber: phoneNumber,
			AccountID:   accountID,
		})); err != nil {
			return err
		}

		if err := tx.CreateEmailAddress(ctx, entities.NewEmailAddress(entities.NewEmailAddressInput{
			Primary:      true,
			EmailAddress: emailAddress,
			AccountID:    accountID,
		})); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &accountV1.CreateAccountResponse{
		Account: &accountV1.Account{
			// TODO populate fields
			Id:             "",
			EmailAddresses: nil,
			PhoneNumbers:   nil,
		},
	}, nil
}
