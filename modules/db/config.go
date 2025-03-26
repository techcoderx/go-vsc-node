package db

import (
	"errors"
	"os"
	"vsc-node/modules/config"
)

var ErrEmptyURI = errors.New("empty MongoDB URI")

type dbConfig struct {
	DbURI string
}

type dbConfigStruct struct {
	*config.Config[dbConfig]
}

type DbConfig = *dbConfigStruct

func NewDbConfig() DbConfig {
	return &dbConfigStruct{config.New(dbConfig{
		DbURI: "mongodb://localhost:27017",
	}, nil)}
}

func (dc *dbConfigStruct) Init() error {
	err := dc.Config.Init()
	if err != nil {
		return err
	}

	url := os.Getenv("MONGO_URL")
	if url != "" {
		return dc.SetDbURI(url)
	}

	return nil
}

func (dc *dbConfigStruct) SetDbURI(uri string) error {
	if uri == "" {
		return ErrEmptyURI
	}
	return dc.Update(func(dc *dbConfig) {
		dc.DbURI = uri
	})
}
