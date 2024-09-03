package db

import "errors"

// Обозначение ошибок
var (
	ErrMigrate       = errors.New("migration failed")
	ErrDuplicate     = errors.New("record already exists")
	ErrNotExist      = errors.New("row does not exist")
	ErrAuthorize     = errors.New("authorize failed")
	ErrYandexSpeller = errors.New("yandex speller failed")
)
