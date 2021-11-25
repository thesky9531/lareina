package sql

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type Config struct {
	Addr        string
	Active      int
	Idle        int
	IdleTimeout int
}

func NewPSQL(c *Config) *sql.DB {
	db, err := sql.Open("postgres", c.Addr)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	db.SetConnMaxLifetime(time.Second * time.Duration(c.IdleTimeout))
	return db
}
