package interfaces

import (
	"context"
	"menti/pkg/domain/account"
	"menti/pkg/domain/note"
)

type ServiceUseCase interface {
	Migrate(ctx context.Context) error

	//Account
	Login(ctx context.Context, acc account.Login) (int64, error)
	//Note
	NoteGetAll(ctx context.Context, userId string) ([]note.Note, error)
	NotePost(ctx context.Context, userId string, note note.Note) (int64, error)
}
