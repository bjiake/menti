//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	http "menti/pkg/api"
	"menti/pkg/api/handler"
	"menti/pkg/config"
	"menti/pkg/db"
	"menti/pkg/service"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectToBD, service.NewService, handler.NewHandler, http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
