package cmd

import (
	"fmt"

	"github.com/meton888/draft/address"
	"github.com/meton888/draft/config"
	"github.com/meton888/draft/docker"
	"github.com/meton888/draft/env"
	"github.com/urfave/cli/v2"
)

var UpCommand = &cli.Command{
	Name:  "up",
	Usage: "Bring the cluster up",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		cfg, _ := config.Yaml()

		zkAddr := address.Zookeeper(cfg.Nodes.Master)
		serversAddr := address.Servers(cfg.Nodes.Master)

		for i, node := range cfg.Nodes.Master {
			cli, _ := docker.Client(address.SSH(node.User, node.Address.External, 0))
			cli.NegotiateAPIVersion(ctx)

			// start
			err := zk.Up(ctx, cli, env.Zookeeper{
				MYID:    i + 1,
				SERVERS: serversAddr,
			})
			if err != nil {
				fmt.Println(err.Error())
			}

			err = master.Up(ctx, cli, env.MesosMaster{
				MESOS_HOSTNAME: node.Address.Internal,
				MESOS_IP:       node.Address.Internal,
				MESOS_ZK:       fmt.Sprintf("%s/mesos", zkAddr),
			})
			if err != nil {
				fmt.Println(err.Error())
			}

			err = marathon.Up(ctx, cli, env.Marathon{
				MARATHON_HOSTNAME:      node.Address.Internal,
				MARATHON_HTTPS_ADDRESS: node.Address.Internal,
				MARATHON_HTTP_ADDRESS:  node.Address.Internal,
				MARATHON_MASTER:        fmt.Sprintf("%s/mesos", zkAddr),
				MARATHON_ZK:            fmt.Sprintf("%s/marathon", zkAddr),
			})
			if err != nil {
				fmt.Println(err.Error())
			}

		}

		for _, node := range cfg.Nodes.Slave {
			cli, _ := docker.Client(address.SSH(node.User, node.Address.External, 0))
			cli.NegotiateAPIVersion(ctx)

			// start
			err := slave.Up(ctx, cli, env.MesosSlave{
				MESOS_HOSTNAME: node.Address.Internal,
				MESOS_IP:       node.Address.Internal,
				MESOS_MASTER:   fmt.Sprintf("%s/mesos", zkAddr),
			})
			if err != nil {
				fmt.Println(err.Error())
			}

		}

		return nil
	},
}
