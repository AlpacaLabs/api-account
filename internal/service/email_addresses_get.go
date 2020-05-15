package service

import (
	"context"

	"github.com/AlpacaLabs/api-account/internal/db"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	paginationV1 "github.com/AlpacaLabs/protorepo-pagination-go/alpacalabs/pagination/v1"
)

// GetEmailAddresses retrieves all email addresses in the system.
// Ideally, this function should be locked down and offered for
// admins only.
func (s Service) GetEmailAddresses(ctx context.Context, request accountV1.GetEmailAddressesRequest) (*accountV1.GetEmailAddressesResponse, error) {
	out := &accountV1.GetEmailAddressesResponse{}

	// Validate cursor
	if request.CursorRequest == nil {
		return nil, ErrNilCursorRequest
	}
	cursorRequest := *request.CursorRequest
	if cursorRequest.Count > MaxPageSize {
		cursorRequest.Count = MaxPageSize
	} else if cursorRequest.Count == 0 {
		cursorRequest.Count = DefaultPageSize
	}

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
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
			out.CursorResponse.NextCursor = emailAddresses[count-1].Id
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}
