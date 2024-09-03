package main

import (
	log "github.com/sirupsen/logrus"
	"menti/pkg/config"
	"menti/pkg/di"
)

func main() {
	cfg, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	server, diErr := di.InitializeAPI(cfg.Db)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start(cfg.Host)
	}
}
