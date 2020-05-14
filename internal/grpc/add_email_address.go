package grpc

import (
	"context"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

func (s Server) AddEmailAddress(ctx context.Context, request *accountV1.AddEmailAddressRequest) (*accountV1.AddEmailAddressResponse, error) {
	return s.service.AddEmailAddress(ctx, request)
}
