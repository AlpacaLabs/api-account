package service

import (
	"context"

	"github.com/rs/xid"
)

func getRequesterID(ctx context.Context) string {
	// TODO this is wrong. instead, parse the actual value from ctx
	return xid.New().String()
}
