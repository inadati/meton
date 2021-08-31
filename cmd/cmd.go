package cmd

import (
	"context"

	"github.com/meton888/meton/config"
)

var (
	ctx    = context.Background()
	cfg, _ = config.Yaml()
)
