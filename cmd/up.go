package cmd

import (
	"fmt"

	"github.com/meton888/meton/config"
	"github.com/meton888/meton/container"
	"github.com/meton888/meton/docker"
	"github.com/meton888/meton/endpoint"
	"github.com/meton888/meton/env"
	"github.com/urfave/cli/v2"
)

var UpCommand = &cli.Command{
	Name:  "up",
	Usage: "Bring the cluster up",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		cfg, _ := config.Yaml()

		svrsAddr, zkAddr := endpoint.MasterNode.AddrCollection(cfg.Cluster.Nodes.Master)

		for i, node := range cfg.Cluster.Nodes.Master {
			dockerClient, _ := docker.Client(endpoint.Node.SSH(cfg.Cluster.Owner, node.Address.External, 0))
			dockerClient.NegotiateAPIVersion(ctx)

			// start zookeeper
			err := container.Zookeeper.Up(ctx, dockerClient, env.Zookeeper{
				MYID:    i + 1,
				SERVERS: svrsAddr,
			},)
			if err != nil {
				fmt.Println(err.Error())
			}

			// start mesos master
			err = container.Mesos.Master.Up(ctx, dockerClient, env.MesosMaster{
				MESOS_HOSTNAME: node.Address.Internal,
				MESOS_IP:       node.Address.Internal,
				MESOS_ZK:       fmt.Sprintf("%s/mesos", zkAddr),
			})
			if err != nil {
				fmt.Println(err.Error())
			}

			// start marathon
			err = container.Marathon.Up(ctx, dockerClient, env.Marathon{
				MARATHON_HOSTNAME:      node.Address.Internal,
				MARATHON_HTTPS_ADDRESS: node.Address.Internal,
				MARATHON_HTTP_ADDRESS:  node.Address.Internal,
				MARATHON_MASTER:        fmt.Sprintf("%s/mesos", zkAddr),
				MARATHON_ZK:            fmt.Sprintf("%s/marathon", zkAddr),
			})
			if err != nil {
				fmt.Println(err.Error())
			}

			// start chronos
			err = container.Chronos.Up(ctx, dockerClient, env.Chronos{
				CHRONOS_MASTER:   fmt.Sprintf("%s/mesos", zkAddr),
				CHRONOS_ZK_HOSTS: zkAddr,
			})
			if err != nil {
				fmt.Println(err.Error())
			}

		}

		for _, node := range cfg.Cluster.Nodes.Slave {
			dockerClient, _ := docker.Client(endpoint.Node.SSH(cfg.Cluster.Owner, node.Address.External, 0))
			dockerClient.NegotiateAPIVersion(ctx)

			// start mesos slave
			err := container.Mesos.Slave.Up(ctx, dockerClient, env.MesosSlave{
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
