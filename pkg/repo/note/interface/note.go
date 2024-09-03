package interfaces

import (
	"context"
	"menti/pkg/domain/note"
)

type NoteRepository interface {
	Migrate(ctx context.Context) error
	GetAll(ctx context.Context, userID int64) ([]note.Note, error)
	Post(ctx context.Context, userID int64, newPeople note.Note) (int64, error)
}
