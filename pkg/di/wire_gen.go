package di

import (
	"context"
	log "github.com/sirupsen/logrus"
	http "menti/pkg/api"
	"menti/pkg/api/handler"
	"menti/pkg/config"
	"menti/pkg/db"
	"menti/pkg/repo/account"
	"menti/pkg/repo/note"
	"menti/pkg/service"
)

func InitializeAPI(cfg config.DataBase) (*http.ServerHTTP, error) {
	bd, err := db.ConnectToBD(cfg)
	if err != nil {
		return nil, err
	}
	log.Info("DB connected")
	// Repository
	accountRepository := account.NewAccountDataBase(bd)
	noteRepository := note.NewNoteDataBase(bd)
	// Хардкод пользователя
	_, err = accountRepository.Registration(context.Background())
	if err != nil {
		switch err {
		case db.ErrDuplicate:
			log.Info("User already registered.")
		default:
			log.Error("Error registration", err)
			return nil, err
		}
	}
	//service - logic
	userService := service.NewService(accountRepository, noteRepository)

	// Init Migrate
	err = userService.Migrate(context.Background())
	if err != nil {
		return nil, err
	}

	userHandler := handler.NewHandler(userService)
	log.Info("handler initialized")
	serverHTTP := http.NewServerHTTP(userHandler)
	log.Info("server http initialized", *serverHTTP)
	return serverHTTP, nil
}
