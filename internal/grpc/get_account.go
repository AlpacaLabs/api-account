package grpc

import (
	"context"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

func (s Server) GetAccount(ctx context.Context, request *accountV1.GetAccountRequest) (*accountV1.GetAccountResponse, error) {
	return s.service.GetAccount(ctx, request)
}
