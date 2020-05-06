package service

import (
	"context"
)

const (
	MinUsernameLength = 4
	MaxUsernameLength = 25
)

// TODO only admins can call this endpoint
func (s *Service) GetAccounts(ctx context.Context) {}

func (s *Service) GetAccount(ctx context.Context) {}

// TODO only admins can create
func (s *Service) CreateAccount(ctx context.Context) {}

func (s *Service) UpdateAccount(ctx context.Context) {}

func (s *Service) DeleteAccount(ctx context.Context) {}
