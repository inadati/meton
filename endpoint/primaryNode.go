package endpoint

import (
	"fmt"

	"github.com/meton888/meton/config"
)

type PrimaryNodeRecipe struct{}

var PrimaryNode = &PrimaryNodeRecipe{}

func (r *PrimaryNodeRecipe) SSH(masterNodes []config.Node) string {
	return fmt.Sprintf("%s:22", masterNodes[0].Address.External)
}

func (r *PrimaryNodeRecipe) Marathon(masterNodes []config.Node) string {
	return fmt.Sprintf("%s:8080", masterNodes[0].Address.Internal)
}

func (r *PrimaryNodeRecipe) MesosMaster(masterNodes []config.Node) string {
	return fmt.Sprintf("%s:5050", masterNodes[0].Address.Internal)
}

func (r *PrimaryNodeRecipe) Chronos(masterNodes []config.Node) string {
	return fmt.Sprintf("%s:4400", masterNodes[0].Address.Internal)
}
