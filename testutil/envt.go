package testutil

import (
	"os"

	"github.com/kcmvp/env/internal"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	_ = os.Setenv(internal.ProfileEnvName, internal.TestProfileValue)
}
