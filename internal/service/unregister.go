package service

import (
	"context"

	"github.com/AlpacaLabs/api-account/internal/db"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

func (s Service) UnregisterEmailAddress(ctx context.Context, request *accountV1.UnregisterEmailAddressRequest) (*accountV1.UnregisterEmailAddressResponse, error) {
	requesterID := getRequesterID(ctx)
	emailAddressID := request.EmailAddressId

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		if e, err := tx.GetEmailAddressByID(ctx, emailAddressID); err != nil {
			return err
		} else if e.Primary {
			return ErrUnregisterPrimaryEmailAddress
		} else if e.AccountId != requesterID {
			return ErrUnregisterUnownedEmailAddress
		}

		_, err := tx.DeleteEmailAddress(ctx, emailAddressID)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &accountV1.UnregisterEmailAddressResponse{}, nil
}

func (s Service) UnregisterPhoneNumber(ctx context.Context, request *accountV1.UnregisterPhoneNumberRequest) (*accountV1.UnregisterPhoneNumberResponse, error) {
	requesterID := getRequesterID(ctx)
	phoneNumberID := request.PhoneNumberId

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		if e, err := tx.GetPhoneNumberByID(ctx, phoneNumberID); err != nil {
			return err
		} else if e.AccountId != requesterID {
			return ErrUnregisterUnownedPhoneNumber
		}

		_, err := tx.DeletePhoneNumber(ctx, phoneNumberID)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &accountV1.UnregisterPhoneNumberResponse{}, nil
}
