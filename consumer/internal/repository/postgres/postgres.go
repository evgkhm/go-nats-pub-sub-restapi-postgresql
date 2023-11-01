package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
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

func NewPostgresDB(config *Conf) (*sqlx.DB, error) {
	dbURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Name, config.Pass, config.SSLMode)

	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgresDB - sqlx.Open: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgresDB - db.Ping: %w", err)
	}

	return db, nil
}
