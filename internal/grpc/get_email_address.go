package grpc

import (
	"context"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

func (s Server) GetEmailAddress(ctx context.Context, request *accountV1.GetEmailAddressRequest) (*accountV1.GetEmailAddressResponse, error) {
	return s.service.GetEmailAddress(ctx, request)
}
