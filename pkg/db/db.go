package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
	"menti/pkg/config"
)

// ConnectToBD Подключение к PostgresSql по app.env
func ConnectToBD(cfg config.DataBase) (*sql.DB, error) {
	// Формирование строки подключения из конфига
	addr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PsqlUser, cfg.PsqlPass, cfg.PsqlHost, cfg.PsqlPort, cfg.PsqlDBName)
	psqlInfo := addr

	// Подключение к БД
	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}
