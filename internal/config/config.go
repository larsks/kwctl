package config

import (
	"log/slog"
)

type (
	Config struct {
		Bps     int
		Verbose int
		Vfo     string
		Device  string
		Pretty  bool
		NoCheck bool
	}

	Context struct {
		Config Config
		Logger *slog.Logger
	}
)
