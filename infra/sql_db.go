package infra

import (
	"database/sql"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/kcmvp/env"
)

var (
	dbOnce sync.Once
	_db    *sql.DB
)

func DefaultDB() *sql.DB {
	dbOnce.Do(func() {
		var err error
		driver := env.Profile().MustGet().GetString("datasource.default.driver")
		url := env.Profile().MustGet().GetString("datasource.default.url")
		if strings.TrimSpace(driver) == "" || strings.TrimSpace(url) == "" {
			slog.Warn("database configuration is missing, skipping database initialization")
			return
		}
		_db, err = sql.Open(driver, url)
		if err != nil {
			panic(err)
		}
		_db.SetConnMaxLifetime(time.Minute * 3)
		_db.SetMaxOpenConns(10)
		_db.SetMaxIdleConns(10)
	})
	return _db
}

func DB(name string) *sql.DB {
	dbOnce.Do(func() {
		var err error
		driver := env.Profile().MustGet().GetString("datasource.default.driver")
		url := env.Profile().MustGet().GetString("datasource.default.url")
		if strings.TrimSpace(driver) == "" || strings.TrimSpace(url) == "" {
			slog.Warn("database configuration is missing, skipping database initialization", "name", name)
			return
		}
		_db, err = sql.Open(driver, url)
		if err != nil {
			panic(err)
		}
		_db.SetConnMaxLifetime(time.Minute * 3)
		_db.SetMaxOpenConns(10)
		_db.SetMaxIdleConns(10)
	})
	return _db
}
