package service

import (
	"context"

	"github.com/AlpacaLabs/api-account/internal/db"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	paginationV1 "github.com/AlpacaLabs/protorepo-pagination-go/alpacalabs/pagination/v1"
)

func (s Service) GetAccount(ctx context.Context, request *accountV1.GetAccountRequest) (*accountV1.GetAccountResponse, error) {
	var accountID string
	var emailAddresses []*accountV1.EmailAddress
	var phoneNumbers []*accountV1.PhoneNumber

	err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {

		// Look up the account by email address ID
		if eid := request.GetEmailAddressId(); eid != "" {
			emailAddress, err := tx.GetEmailAddressByID(ctx, eid)
			if err != nil {
				return err
			}

			accountID = emailAddress.AccountId

			e, err := tx.GetEmailAddressesForAccount(ctx, accountID, paginationV1.CursorRequest{
				// A user can't have more than 20 email addresses, right??
				Count: 20,
			})
			if err != nil {
				return err
			}

			emailAddresses = e
		} else if pid := request.GetPhoneNumberId(); pid != "" {
			phoneNumber, err := tx.GetPhoneNumberByID(ctx, pid)
			if err != nil {
				return err
			}

			accountID = phoneNumber.AccountId

			p, err := tx.GetPhoneNumbersForAccount(ctx, accountID, paginationV1.CursorRequest{
				// A user can't have more than 20 email addresses, right??
				Count: 20,
			})
			if err != nil {
				return err
			}

			phoneNumbers = p
		}

		// TODO support other account identifiers

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &accountV1.GetAccountResponse{
		Account: &accountV1.Account{
			Id:             accountID,
			EmailAddresses: emailAddresses,
			PhoneNumbers:   phoneNumbers,
		},
	}, nil

}
