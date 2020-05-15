package grpc

import (
	"context"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

func (s Server) RegisterEmailAddress(ctx context.Context, request *accountV1.RegisterEmailAddressRequest) (*accountV1.RegisterEmailAddressResponse, error) {
	return s.service.RegisterEmailAddress(ctx, request)
}

func (s Server) RegisterPhoneNumber(ctx context.Context, request *accountV1.RegisterPhoneNumberRequest) (*accountV1.RegisterPhoneNumberResponse, error) {
	return s.service.RegisterPhoneNumber(ctx, request)
}
