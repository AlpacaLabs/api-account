package service

import (
	"context"

	"github.com/AlpacaLabs/api-account/internal/db"
	authV1 "github.com/AlpacaLabs/protorepo-auth-go/alpacalabs/auth/v1"
	paginationV1 "github.com/AlpacaLabs/protorepo-pagination-go/alpacalabs/pagination/v1"
)

const (
	DefaultPageSize = 5
	MaxPageSize     = 1000
)

// GetEmailAddresses retrieves all email addresses in the system.
// Ideally, this function should be locked down and offered for
// admins only.
func (s *Service) GetEmailAddresses(ctx context.Context, request authV1.GetEmailAddressesRequest) (*authV1.GetEmailAddressesResponse, error) {
	out := &authV1.GetEmailAddressesResponse{}

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
func (s *Service) CreateEmailAddress(ctx context.Context) {
	// Validate input. Validate email with checkmail library.
	// Check if email already exists. if yes, resend confirmation email. if no, return error.
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
