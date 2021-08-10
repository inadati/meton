package cmd

import (
	"fmt"

	"github.com/meton888/draft/address"
	"github.com/meton888/draft/config"
	"github.com/meton888/draft/container"
	"github.com/meton888/draft/docker"
	"github.com/urfave/cli/v2"
)

var DestroyCommand = &cli.Command{
	Name:  "destroy",
	Usage: "Teardown the cluster and clean cluster nodes",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		cfg, _ := config.Yaml()

		for _, node := range cfg.Nodes.Master {
			cli, _ := docker.Client(address.SSH(node.User, node.Address.External, 0))
			cli.NegotiateAPIVersion(ctx)

			// stop
			err := container.Down(ctx, cli)
			if err != nil {
				fmt.Println(err.Error())
			}

		}

		for _, node := range cfg.Nodes.Slave {
			cli, _ := docker.Client(address.SSH(node.User, node.Address.External, 0))
			cli.NegotiateAPIVersion(ctx)

			// stop
			err := container.Down(ctx, cli)
			if err != nil {
				fmt.Println(err.Error())
			}

		}
		return nil
	},
}
