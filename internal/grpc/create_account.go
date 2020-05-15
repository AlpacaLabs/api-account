package grpc

import (
	"context"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
)

func (s Server) CreateAccount(ctx context.Context, request *accountV1.CreateAccountRequest) (*accountV1.CreateAccountResponse, error) {
	return s.service.CreateAccount(ctx, request)
}
