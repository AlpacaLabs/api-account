package grpc

import (
	"context"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Server) AddPhoneNumber(ctx context.Context, request *accountV1.AddPhoneNumberRequest) (*accountV1.AddPhoneNumberResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}