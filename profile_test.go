package env

import (
	"testing"

	_ "github.com/kcmvp/env/testutil"
	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	url := Profile().MustGet().GetString("datasource.default.url")
	driver := Profile().MustGet().GetString("datasource.default.driver")
	assert.Equal(t, "file::memory:?cache=shared", url)
	assert.Equal(t, "sqlite3", driver)

	url = Profile().MustGet().GetString("datasource.slave.url")
	driver = Profile().MustGet().GetString("datasource.slave.driver")
	assert.Equal(t, "192.168", url)
	assert.Equal(t, "mysql", driver)
}
