package config

import (
	"log/slog"
)

type (
	Config struct {
		Bitrate int
		Verbose int
		Vfo     string
		Device  string
	}

	Context struct {
		Config Config
		Logger *slog.Logger
	}
)
