package db

import (
	"demo-1/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDb(conf *configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(conf.Dsn.DsnName), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Db{db}
}
