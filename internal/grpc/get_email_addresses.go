package grpc

import (
	"context"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

func (s Server) GetEmailAddresses(ctx context.Context, request *accountV1.GetEmailAddressesRequest) (*accountV1.GetEmailAddressesResponse, error) {
	return s.service.GetEmailAddresses(ctx, request)
}
