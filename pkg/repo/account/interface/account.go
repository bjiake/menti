package interfaces

import (
	"context"
	"menti/pkg/domain/account"
)

type AccountRepository interface {
	Migrate(ctx context.Context) error
	Registration(ctx context.Context) (*account.Account, error)
	Login(ctx context.Context, acc account.Login) (int64, error)
	CheckAccount(ctx context.Context, id int64) error
}
