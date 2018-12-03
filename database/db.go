package database

import (
	"github.com/finalist736/gokit/logger"
	"github.com/gocraft/dbr"
	"github.com/gocraft/health"
	"time"
)

type DBConfig struct {
	Driver string
	// username:password@tcp(host:port)/dbname
	Dsn            string
	ConnectionName string
	MaxOpenConns   int
	MaxIdleConns   int
	LifeTime       time.Duration
	Stream         *health.Stream
}

var connections map[string]*dbr.Connection

func init() {
	connections = make(map[string]*dbr.Connection, 1)
}

func Add(cfg *DBConfig) error {
	if cfg.Stream == nil {
		cfg.Stream = logger.DatabaseStream()
	}
	conn, err := dbr.Open(cfg.Driver, cfg.Dsn, cfg.Stream)
	if err != nil {
		return err
	}
	conn.SetMaxOpenConns(cfg.MaxOpenConns)
	conn.SetMaxIdleConns(cfg.MaxIdleConns)
	conn.SetConnMaxLifetime(cfg.LifeTime)
	err = conn.Ping()
	if err != nil {
		return err
	}
	if cfg.ConnectionName == "" {
		cfg.ConnectionName = "default"
	}
	connections[cfg.ConnectionName] = conn
	return nil
}

func GetConnection(name string) *dbr.Connection {
	return connections[name]
}

func GetSession(name string) *dbr.Session {
	c, ok := connections[name]
	if !ok {
		return nil
	}
	return c.NewSession(nil)
}

func GetDefaultSession() *dbr.Session {
	return GetSession("default")
}
