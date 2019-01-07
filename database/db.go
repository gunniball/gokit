package database

import (
	"errors"
	"sync"
	"time"

	"github.com/finalist736/gokit/logger"
	"github.com/gocraft/dbr"
	"github.com/gocraft/health"
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

var mux sync.RWMutex
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
	mux.Lock()
	connections[cfg.ConnectionName] = conn
	mux.Unlock()
	return nil
}

func Remove(name string) error {
	conn := GetConnection(name)
	if conn == nil {
		return errors.New("no such connection")
	}
	err := conn.Close()
	if err != nil {
		return nil
	}
	mux.Lock()
	delete(connections, name)
	mux.Unlock()
	return nil
}

func Close() {
	mux.Lock()
	var err error
	for _, iterator := range connections {
		err = iterator.Close()
		if err != nil {
			logger.StdErr().Warnf("Close connection error: %s", err)
		}
	}
	connections = make(map[string]*dbr.Connection, 0)
	mux.Unlock()
}

func GetConnection(name string) *dbr.Connection {
	mux.RLock()
	defer mux.RUnlock()
	return connections[name]
}

func GetSession(name string) *dbr.Session {
	mux.RLock()
	c, ok := connections[name]
	mux.RUnlock()
	if !ok {
		return nil
	}
	return c.NewSession(nil)
}

func GetDefaultSession() *dbr.Session {
	return GetSession("default")
}
