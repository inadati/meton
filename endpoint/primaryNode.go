package endpoint

import (
	"fmt"

	"github.com/meton888/meton/config"
)

type primaryNode struct {
	SSH         func(masterNodes []config.Node) string
	Marathon    func(masterNodes []config.Node) string
	MesosMaster func(masterNodes []config.Node) string
	Chronos     func(masterNodes []config.Node) string
}

var PrimaryNode = &primaryNode{
	SSH: func(masterNodes []config.Node) string {
		return fmt.Sprintf("%s:22", masterNodes[0].Address.External)
	},
	Marathon: func(masterNodes []config.Node) string {
		return fmt.Sprintf("%s:8080", masterNodes[0].Address.Internal)
	},
	MesosMaster: func(masterNodes []config.Node) string {
		return fmt.Sprintf("%s:5050", masterNodes[0].Address.Internal)
	},
	Chronos: func(masterNodes []config.Node) string {
		return fmt.Sprintf("%s:4400", masterNodes[0].Address.Internal)
	},
}