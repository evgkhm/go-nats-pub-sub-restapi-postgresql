package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log/slog"
)

var Config *Conf

type Conf struct {
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	SSLMode string
}

func (c Conf) Init() {
	Config = &Conf{
		Host:    viper.GetString("postgres.host"),
		Port:    viper.GetString("postgres.port"),
		User:    viper.GetString("postgres.user"),
		Name:    viper.GetString("postgres.name"),
		Pass:    viper.GetString("postgres.pass"),
		SSLMode: viper.GetString("postgres.ssl_mode"),
	}
}

func NewPostgresDB(config *Conf, logger *slog.Logger) (*sqlx.DB, error) {
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Name, config.Pass, config.SSLMode)

	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		logger.Error("postgres - NewPostgresDB - sqlx.Open", "err", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error("postgres - NewPostgresDB - db.Ping", "err", err)
		return nil, err
	}

	return db, nil
}
