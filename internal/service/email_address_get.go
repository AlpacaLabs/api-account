package service

import (
	"context"

	"github.com/AlpacaLabs/api-account/internal/db"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

// GetEmailAddress retrieves an email address by primary key.
func (s Service) GetEmailAddress(ctx context.Context, request *accountV1.GetEmailAddressRequest) (*accountV1.GetEmailAddressResponse, error) {
	requesterID := getRequesterID(ctx)
	response := &accountV1.GetEmailAddressResponse{}
	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		emailAddress, err := tx.GetEmailAddressByID(ctx, request.Id)
		if err != nil {
			return err
		}

		if emailAddress.AccountId != requesterID {
			return ErrUnowned
		}

		response.EmailAddress = emailAddress
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
