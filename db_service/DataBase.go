package dbservice

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"Medods/utils"

	_ "github.com/lib/pq"
)

type Config struct {
	DataBase DataBaseConfig `json:"Database"`
}

type DataBaseConfig struct {
	Port     int    `json:"port"`
	Host     string `json:"host"`
	User     string `json:"user"`
	DbName   string `json:"db_name"`
	Password string `json:"password"`
	SSLMode  string `json:"sslmode"`
}

func Connect() (*sql.DB, error) {
	var config Config

	data, err := utils.FileByName("config.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	dbConfig := config.DataBase

	dataSource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DbName)

	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	schema, err := utils.FileByName("scheme.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
