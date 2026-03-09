package db

import (
	"database/sql"
	"testing"

	_ "github.com/kcmvp/env/testutil"
	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	assert.ElementsMatch(t, []string{"default", "slave"}, All())
	dbs.Range(func(key, value interface{}) bool {
		name := key.(string)
		db := value.(*sql.DB)
		assert.NotNil(t, db, "db should not be nil for datasource %s", name)
		return true
	})
}
