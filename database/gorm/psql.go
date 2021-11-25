package gorm

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type Config struct {
	Addr        string
	Active      int
	Idle        int
	IdleTimeout int
}

func NewPSQL(c *Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(c.Addr), &gorm.Config{
		Logger: logger.New(logrus.StandardLogger(), logger.Config{}),
	})
	if err != nil {
		panic(err)
	}
	sdb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sdb.SetMaxIdleConns(c.Idle)
	sdb.SetConnMaxLifetime(time.Second * time.Duration(c.IdleTimeout))
	return db
}
