package env

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spf13/viper"
)

const (
	envProfile  = "profile"
	testProfile = "test"
	cfgName     = "application"
	testCfgName = "application_test"
)

var (
	_cfg *viper.Viper
	once sync.Once
)

// Profile loads the application configuration.
//
// Rules:
//  1. If the current process is running `go test`, it tries application_test.yml.
//  2. Otherwise it tries application.yml.
//  3. It searches in '.' and './config'.
//
// If required==false, missing config is not an error.
func Profile() mo.Result[*viper.Viper] {
	once.Do(func() {
		_cfg, _ = loadViper(false)
	})
	return lo.If(_cfg == nil, mo.Err[*viper.Viper](fmt.Errorf("can not find application.yml"))).Else(mo.Ok(_cfg))
}

func loadViper(required bool) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Search paths:
	// 1) project root (detected by walking up to find go.mod) and its ./config
	// 2) current working directory and its ./config
	//
	// This makes discovery work for both:
	// - dev time: running from repo root, package dirs, or IDE
	// - runtime: running from an app's working directory
	addDefaultConfigPaths(v)
	name := cfgName
	if os.Getenv(envProfile) == testProfile {
		name = testCfgName
	}
	v.SetConfigName(strings.TrimSuffix(name, filepath.Ext(name)))
	// previous test-specific explicit checks removed
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !required && errors.As(err, &notFound) {
			return v, nil
		}
		return nil, fmt.Errorf("read %s: %w", name, err)
	}
	return v, nil
}

// addDefaultConfigPaths registers a stable set of config search paths into viper.
//
// Why this exists:
// Viper resolves relative paths against the *current working directory* (CWD). In dev-time (IDE,
// `go test`, running from package folders), CWD can vary a lot. In runtime, users often expect
// config to be read from the directory they launch the binary from.
//
// Strategy (in order):
//  1. Project root (nearest parent dir containing go.mod) and its "config" subdir.
//     - This makes dev-time stable: regardless of running tests from ./app, ./cmd/..., etc,
//     we can still find config at the repo root.
//  2. Current working directory and its "config" subdir.
//     - This preserves runtime flexibility: users can ship config next to the binary or run from
//     a deployment directory that contains application.yml.
//
// Notes:
//   - If a path is added multiple times, viper will just search it multiple times; that's harmless.
//   - If CWD can't be determined, we fall back to "." and "./config".
func addDefaultConfigPaths(v *viper.Viper) {
	cwd, err := os.Getwd()
	if err != nil {
		// Worst-case fallback.
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		return
	}

	if root, ok := findProjectRoot(cwd); ok {
		v.AddConfigPath(root)
		v.AddConfigPath(filepath.Join(root, "config"))
	}

	// Always fall back to CWD.
	v.AddConfigPath(cwd)
	v.AddConfigPath(filepath.Join(cwd, "config"))
}

// findProjectRoot walks upward from `start` until it finds a directory containing a go.mod.
// It returns (root, true) if found, otherwise ("", false).
//
// This is intentionally lightweight and filesystem-only:
// we don't parse go.mod; the existence check is sufficient for locating the module root.
func findProjectRoot(start string) (string, bool) {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
