package service

import (
	"context"
	"menti/pkg/domain/account"
)

func (s *service) Login(ctx context.Context, acc account.Login) (int64, error) {
	id, err := s.rAccount.Login(ctx, acc)
	if err != nil {
		return 0, err
	}
	return id, nil
}
