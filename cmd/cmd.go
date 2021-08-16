package cmd

import (
	"context"

	"github.com/meton888/meton/address"
	"github.com/meton888/meton/container"
)

var (
	ctx          = context.Background()
	zk           container.Zookeeper
	master       container.MesosMaster
	marathon     container.Marathon
	chronos      container.Chronos
	slave        container.MesosSlave
	compoundAddr address.Compound
)
