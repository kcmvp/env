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
	dbCfgName     = "datasource"
	defaultDBName = "default"
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
		lo.ForEach(lo.Keys(dsm), func(ds string, index int) {
			key := fmt.Sprintf("datasource.%s.driver", ds)
			driver := cfg.GetString(key)
			key = fmt.Sprintf("datasource.%s.url", ds)
			url := cfg.GetString(key)
			if strings.TrimSpace(driver) == "" || strings.TrimSpace(url) == "" {
				panic(fmt.Sprintf("driver or url is missing for database %s", ds))
			}
			_db, err := sql.Open(driver, url)
			if err != nil {
				panic(fmt.Sprintf("open database %s failed: %v", ds, err))
			}
			_db.SetConnMaxLifetime(time.Minute * 3)
			_db.SetMaxOpenConns(10)
			_db.SetMaxIdleConns(10)
			dbs.Store(ds, _db)
		})
	})
	v, ok := dbs.Load(name)
	if !ok {
		return nil
	}
	return v.(*sql.DB)
}
