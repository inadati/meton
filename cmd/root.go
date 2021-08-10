package cmd

import (
	"context"
	"github.com/meton888/draft/container"
)

var (
	ctx      = context.Background()
	zk       container.Zookeeper
	master   container.MesosMaster
	marathon container.Marathon
	slave    container.MesosSlave
)
