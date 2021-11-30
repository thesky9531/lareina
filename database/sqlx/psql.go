package sqlx

import (
	_ "database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type Config struct {
	Addr        string `toml:"Addr"`
	Active      int
	Idle        int
	IdleTimeout int
}

func NewPSQL(c *Config) *sqlx.DB {
	db, err := sqlx.Connect("postgres", c.Addr)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	db.SetConnMaxLifetime(time.Second * time.Duration(c.IdleTimeout))
	return db
}
