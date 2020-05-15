package service

import (
	"context"
	"fmt"

	"github.com/AlpacaLabs/api-account/internal/db"
	"github.com/AlpacaLabs/api-account/internal/db/entities"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	paginationV1 "github.com/AlpacaLabs/protorepo-pagination-go/alpacalabs/pagination/v1"
	"github.com/badoux/checkmail"
	"github.com/ttacon/libphonenumber"
)

// RegisterEmailAddress creates an email address entity for a
// given email address and account ID.
func (s Service) RegisterEmailAddress(ctx context.Context, request *accountV1.RegisterEmailAddressRequest) (*accountV1.RegisterEmailAddressResponse, error) {
	emailAddress := request.EmailAddress
	accountID := request.AccountId

	// Validate email address format
	if err := checkmail.ValidateFormat(emailAddress); err != nil {
		return nil, err
	}

	out := &accountV1.RegisterEmailAddressResponse{}

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {

		// Does the account already exist?
		// TODO ideally this should come from request context rather than a protobuf field.
		//  clients shouldn't be allowed to register email addresses for other users.
		if _, err := tx.GetAccountByID(ctx, accountID); err != nil {
			return fmt.Errorf("no account found for id: %s", accountID)
		}

		// Is the email already registered?
		email, err := tx.GetEmailAddressByEmailAddress(ctx, emailAddress)

		// Check for internal errors
		if err != nil && err != db.ErrNotFound {
			return err
		}

		// If the email isn't found, then no one has registered it.
		if err == db.ErrNotFound || email == nil {

			// Get their email addresses to determine if the one they want to register
			// should be automatically set to primary
			var isFirstEmailRegistered bool
			if emailAddresses, err := tx.GetEmailAddressesForAccount(ctx, accountID, paginationV1.CursorRequest{
				// A user can't have more than 20 email addresses, right??
				Count: 20,
			}); err != nil {
				if err == db.ErrNotFound {
					isFirstEmailRegistered = true
				} else {
					return err
				}
			} else if len(emailAddresses) == 0 {
				isFirstEmailRegistered = true
			}

			// Create an email address record
			if err := tx.CreateEmailAddress(ctx, entities.NewEmailAddress(entities.NewEmailAddressInput{
				Primary:      isFirstEmailRegistered,
				EmailAddress: emailAddress,
				AccountID:    accountID,
			})); err != nil {
				return err
			}
		} else {
			if email.AccountId != accountID {
				return ErrEmailAlreadyRegisteredByDifferentAccount
			} else {
				// TODO add transactional outbox record to resend confirmation email
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (s Service) RegisterPhoneNumber(ctx context.Context, request *accountV1.RegisterPhoneNumberRequest) (*accountV1.RegisterPhoneNumberResponse, error) {
	phoneNumber := request.PhoneNumber
	accountID := request.AccountId

	// Validate email address format
	if _, err := libphonenumber.Parse(phoneNumber, "US"); err != nil {
		return nil, err
	}

	out := &accountV1.RegisterPhoneNumberResponse{}

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {

		// Does the account already exist?
		// TODO ideally this should come from request context rather than a protobuf field.
		//  clients shouldn't be allowed to register phone numbers for other users.
		if _, err := tx.GetAccountByID(ctx, accountID); err != nil {
			return fmt.Errorf("no account found for id: %s", accountID)
		}

		// Is the phone number already registered?
		entity, err := tx.GetPhoneNumberByPhoneNumber(ctx, phoneNumber)

		// Check for internal errors
		if err != nil && err != db.ErrNotFound {
			return err
		}

		// If the phoneNumber isn't found, then no one has registered it.
		if err == db.ErrNotFound || entity == nil {

			// Create a phone number record
			if err := tx.CreatePhoneNumber(ctx, entities.NewPhoneNumber(entities.NewPhoneNumberInput{
				PhoneNumber: phoneNumber,
				AccountID:   accountID,
			})); err != nil {
				return err
			}
		} else {
			if entity.AccountId != accountID {
				return ErrPhoneAlreadyRegisteredByDifferentAccount
			} else {
				// TODO add transactional outbox record to resend confirmation SMS
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}
