package config

import (
	"github.com/hs-zavet/comtools/logkit"
)

func (c *Config) LogSetup() {
	logger := logkit.SetupLogger(c.Server.Log.Level, c.Server.Log.Format)
	c.Log = logger
}
