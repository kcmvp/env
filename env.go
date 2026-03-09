package env

import (
	"github.com/kcmvp/env/internal"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func Profile() *viper.Viper {
	p := internal.Profile()
	return lo.IfF(p.IsOk(), p.MustGet).Else(nil)
}
