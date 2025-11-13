package config

import (
	"log/slog"
)

type (
	Config struct {
		Bitrate string
		Verbose int
		Vfo     string
		Device  string
	}

	Context struct {
		Config Config
		Logger *slog.Logger
	}
)
