package testutil

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	_ = os.Setenv("profile", "test")
}
