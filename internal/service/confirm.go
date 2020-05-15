package service

import (
	"context"

	"github.com/AlpacaLabs/api-account/internal/db"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

func (s Service) ConfirmEmailAddress(ctx context.Context, request *accountV1.ConfirmEmailAddressRequest) error {
	return s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		return tx.ConfirmEmailAddress(ctx, request.EmailAddressId)
	})
}

func (s Service) ConfirmPhoneNumber(ctx context.Context, request *accountV1.ConfirmPhoneNumberRequest) error {
	return s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		return tx.ConfirmPhoneNumber(ctx, request.PhoneNumberId)
	})
}
