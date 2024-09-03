package service

import (
	"context"
	"menti/pkg/domain/note"
)

func (s *service) NoteGetAll(ctx context.Context, userId string) ([]note.Note, error) {
	id, err := s.checkLogin(ctx, userId)
	if err != nil {
		return nil, err
	}

	notes, err := s.rNote.GetAll(ctx, id)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *service) NotePost(ctx context.Context, userId string, note note.Note) (int64, error) {
	id, err := s.checkLogin(ctx, userId)
	if err != nil {
		return 0, err
	}
	err = s.checkText(note)
	if err != nil {
		return 0, err
	}

	result, err := s.rNote.Post(ctx, id, note)
	if err != nil {
		return 0, err
	}
	return result, nil
}
