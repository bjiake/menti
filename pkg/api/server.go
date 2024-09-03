package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"menti/pkg/api/handler"
	"menti/pkg/config"
	"net/http"
	"time"
)

type ServerHTTP struct {
	router *chi.Mux
}

/*
Договора
*/
func NewServerHTTP(userHandler *handler.Handler) *ServerHTTP {
	// TODO: Авторизация кривая, почему-то не работает, почему??? Всегда работает
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(15 * time.Second))

	router.Get("/login", userHandler.Login)
	router.Route("/notes", userHandler.Note)

	log.Info("Server ready to start")
	return &ServerHTTP{router: router}
}

func (sh *ServerHTTP) Start(cfg config.Server) {
	log.Info("Server starting")
	addr := cfg.IP + ":" + cfg.Port
	err := http.ListenAndServe(addr, sh.router)
	if err != nil {
		log.Fatal(err)
		return
	}
}
