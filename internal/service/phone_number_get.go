package service

import (
	"context"

	"github.com/AlpacaLabs/api-account/internal/db"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

// GetPhoneNumber retrieves an phone number by primary key.
func (s Service) GetPhoneNumber(ctx context.Context, request *accountV1.GetPhoneNumberRequest) (*accountV1.GetPhoneNumberResponse, error) {
	requesterID := getRequesterID(ctx)
	response := &accountV1.GetPhoneNumberResponse{}
	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		phoneNumber, err := tx.GetPhoneNumberByID(ctx, request.Id)
		if err != nil {
			return err
		}

		if phoneNumber.AccountId != requesterID {
			return ErrUnowned
		}

		response.PhoneNumber = phoneNumber
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
