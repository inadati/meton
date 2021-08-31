package cmd

import (
	"fmt"

	"github.com/meton888/meton/container"
	"github.com/meton888/meton/docker"
	"github.com/meton888/meton/endpoint"
	"github.com/urfave/cli/v2"
)

var DownCommand = &cli.Command{
	Name:  "down",
	Usage: "Teardown the cluster and clean cluster nodes",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "target",
			Aliases: []string{"t"},
			Usage:   "Specifies the address.external of the node to be destroyed.",
		},
	},
	Action: func(c *cli.Context) error {
		targetNode := c.String("target")

		if targetNode != "" {
			dockerClient, _ := docker.Client(endpoint.Node.SSH(cfg.Cluster.Owner, targetNode, 0))
			dockerClient.NegotiateAPIVersion(ctx)

			// stop
			err := container.All.Down(ctx, dockerClient)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {

			for _, node := range cfg.Cluster.Nodes.Master {
				dockerClient, _ := docker.Client(endpoint.Node.SSH(cfg.Cluster.Owner, node.Address.External, 0))
				dockerClient.NegotiateAPIVersion(ctx)

				// stop
				err := container.All.Down(ctx, dockerClient)
				if err != nil {
					fmt.Println(err.Error())
				}

			}

			for _, node := range cfg.Cluster.Nodes.Slave {
				dockerClient, _ := docker.Client(endpoint.Node.SSH(cfg.Cluster.Owner, node.Address.External, 0))
				dockerClient.NegotiateAPIVersion(ctx)

				// stop
				err := container.All.Down(ctx, dockerClient)
				if err != nil {
					fmt.Println(err.Error())
				}

			}

		}

		return nil
	},
}
