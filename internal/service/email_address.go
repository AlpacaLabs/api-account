package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlpacaLabs/api-account/internal/db/entities"
	"github.com/badoux/checkmail"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"

	"github.com/AlpacaLabs/api-account/internal/db"
	paginationV1 "github.com/AlpacaLabs/protorepo-pagination-go/alpacalabs/pagination/v1"
)

var (
	ErrEmailAlreadyRegisteredByDifferentAccount = errors.New("that email address is already by a different account")
)

const (
	DefaultPageSize = 5
	MaxPageSize     = 1000
)

// GetEmailAddresses retrieves all email addresses in the system.
// Ideally, this function should be locked down and offered for
// admins only.
func (s *Service) GetEmailAddresses(ctx context.Context, request accountV1.GetEmailAddressesRequest) (*accountV1.GetEmailAddressesResponse, error) {
	out := &accountV1.GetEmailAddressesResponse{}

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		cursorRequest := *request.CursorRequest
		emailAddresses, err := tx.GetEmailAddresses(ctx, cursorRequest)
		if err != nil {
			return err
		}

		out.EmailAddresses = emailAddresses

		count := len(emailAddresses)

		out.CursorResponse = &paginationV1.CursorResponse{
			PreviousCursor: cursorRequest.Cursor,
			Count:          int32(count),
		}

		if count > 0 {
			// Set NextCursor so clients can continue pagination
			out.CursorResponse.NextCursor = emailAddresses[len(emailAddresses)-1].Id
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}

// GetEmailAddress retrieves an email address by primary key.
// This function should only return email addresses that belong
// to you.
func (s *Service) GetEmailAddress(ctx context.Context) {
	// Look up email by ID
}

// CreateEmailAddress creates an email address entity for a
// given email address and account ID.
func (s *Service) AddEmailAddress(ctx context.Context, request *accountV1.AddEmailAddressRequest) (*accountV1.AddEmailAddressResponse, error) {
	emailAddress := request.EmailAddress
	accountID := request.AccountId

	// Validate email address format
	if err := checkmail.ValidateFormat(emailAddress); err != nil {
		return nil, err
	}

	out := &accountV1.AddEmailAddressResponse{}

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

// UpdateEmailAddress updates the email address's confirmation status.
// This is usually done when a user clicks the confirmation link
// in an email they receive.
func (s *Service) UpdateEmailAddress(ctx context.Context) {
	// Check if entity exists for email address
	// If not, return NotFound
	// Update the email's confirmation status
	// Return new entity in response
}

func (s *Service) DeleteEmailAddress(ctx context.Context) {
	// TODO check existence
	// Delete email address
}
