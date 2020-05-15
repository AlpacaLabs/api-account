package grpc

import (
	"context"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

func (s Server) UnregisterEmailAddress(ctx context.Context, request *accountV1.UnregisterEmailAddressRequest) (*accountV1.UnregisterEmailAddressResponse, error) {
	return s.service.UnregisterEmailAddress(ctx, request)
}

func (s Server) UnregisterPhoneNumber(ctx context.Context, request *accountV1.UnregisterPhoneNumberRequest) (*accountV1.UnregisterPhoneNumberResponse, error) {
	return s.service.UnregisterPhoneNumber(ctx, request)
}
