package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kcmvp/env/internal"
	"github.com/samber/lo"
)

var (
	dbOnce sync.Once
	dbs    sync.Map
)

const (
	dbCfgName                    = "datasource"
	defaultDBName                = "default"
	driver                       = "driver"
	url                          = "url"
	maxIdeal                     = "max_idle"
	maxOpen                      = "max_open"
	maxLifeTime                  = "max_lifetime"
	DefMaxOpen                   = 50
	DefMaxIdeal                  = 10
	DefMaxLifeTime time.Duration = time.Minute * 30
)

func Default() *sql.DB {
	return DB(defaultDBName)
}

func All() []string {
	_ = DB(defaultDBName) // trigger dbOnce to load all dbs
	var names []string
	dbs.Range(func(key, value interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}

func DB(name string) *sql.DB {
	dbOnce.Do(func() {
		cfg := internal.MstProfile()
		dsm := cfg.GetStringMap(dbCfgName)
		if dsm == nil {
			panic("database config not found")
		}
		lo.ForEach(lo.Keys(dsm), func(ds string, _ int) {
			_driver := cfg.GetString(fmt.Sprintf("datasource.%s.%s", ds, driver))
			_url := cfg.GetString(fmt.Sprintf("datasource.%s.%s", ds, url))
			if strings.TrimSpace(_driver) == "" || strings.TrimSpace(_url) == "" {
				panic(fmt.Sprintf("driver or url is missing for database %s", ds))
			}
			_db, err := sql.Open(_driver, _url)
			if err != nil {
				panic(fmt.Sprintf("open database %s failed: %v", ds, err))
			}

			_open := cfg.GetInt(fmt.Sprintf("datasource.%s.%s", ds, maxOpen))
			if _open == 0 {
				_open = DefMaxOpen
			}
			_idle := cfg.GetInt(fmt.Sprintf("datasource.%s.%s", ds, maxIdeal))
			if _idle == 0 {
				_idle = DefMaxIdeal
			}
			_maxLifeTime := cfg.GetInt(fmt.Sprintf("datasource.%s.%s", ds, maxLifeTime))
			if _maxLifeTime == 0 {
				_db.SetConnMaxLifetime(DefMaxLifeTime)
			} else {
				_db.SetConnMaxLifetime(time.Minute * time.Duration(_maxLifeTime))
			}
			_db.SetMaxOpenConns(_open)
			_db.SetMaxIdleConns(_idle)
			dbs.Store(ds, _db)
		})
	})
	v, ok := dbs.Load(name)
	if !ok {
		return nil
	}
	return v.(*sql.DB)
}
