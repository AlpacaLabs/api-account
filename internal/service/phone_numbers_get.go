package service

import (
	"context"

	"github.com/AlpacaLabs/api-account/internal/db"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	paginationV1 "github.com/AlpacaLabs/protorepo-pagination-go/alpacalabs/pagination/v1"
)

// GetPhoneNumbers retrieves all phone numbers in the system.
// Ideally, this function should be locked down and offered for
// admins only.
func (s Service) GetPhoneNumbers(ctx context.Context, request *accountV1.GetPhoneNumbersRequest) (*accountV1.GetPhoneNumbersResponse, error) {
	out := &accountV1.GetPhoneNumbersResponse{}

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
		phoneNumbers, err := tx.GetPhoneNumbers(ctx, cursorRequest)
		if err != nil {
			return err
		}

		out.PhoneNumbers = phoneNumbers

		count := len(phoneNumbers)

		out.CursorResponse = &paginationV1.CursorResponse{
			PreviousCursor: cursorRequest.Cursor,
			Count:          int32(count),
		}

		if count > 0 {
			// Set the next cursor clients should use so they can continue pagination
			out.CursorResponse.NextCursor = phoneNumbers[count-1].Id
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return out, nil
}
