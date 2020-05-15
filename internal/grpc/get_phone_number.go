package grpc

import (
	"context"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

func (s Server) GetPhoneNumber(ctx context.Context, request *accountV1.GetPhoneNumberRequest) (*accountV1.GetPhoneNumberResponse, error) {
	return s.service.GetPhoneNumber(ctx, request)
}
