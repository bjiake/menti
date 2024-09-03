package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	Db   DataBase
	Host Server
}
type DataBase struct {
	PsqlUser   string
	PsqlPass   string
	PsqlHost   string
	PsqlPort   string
	PsqlDBName string
}
type Server struct {
	IP   string
	Port string
}

func LoadConfig() (Config, error) {
	var config Config

	//Загрузка переменных окружения (если используете godotenv)
	err := godotenv.Load("./app.env") // Раскомментируйте эту строку, если используете godotenv
	if err != nil {
		log.Fatal("Error loading .env file")
		return config, err
	}
	// Получение значений переменных окружения
	config.Db.PsqlUser = os.Getenv("POSTGRES_USER")
	config.Db.PsqlPass = os.Getenv("POSTGRES_PASSWORD")
	config.Db.PsqlHost = os.Getenv("POSTGRES_HOST")
	config.Db.PsqlDBName = os.Getenv("POSTGRES_DB")
	config.Host.IP = os.Getenv("IP")
	config.Host.Port = os.Getenv("PORT")
	return config, err
}
