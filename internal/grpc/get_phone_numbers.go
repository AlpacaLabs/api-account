package grpc

import (
	"context"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

func (s Server) GetPhoneNumbers(ctx context.Context, request *accountV1.GetPhoneNumbersRequest) (*accountV1.GetPhoneNumbersResponse, error) {
	return s.service.GetPhoneNumbers(ctx, request)
}
