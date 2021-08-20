package cmd

import (
	"fmt"

	"github.com/meton888/meton/address"
	"github.com/meton888/meton/config"
	"github.com/meton888/meton/container"
	"github.com/meton888/meton/docker"
	"github.com/meton888/meton/env"
	"github.com/urfave/cli/v2"
)

var UpCommand = &cli.Command{
	Name:  "up",
	Usage: "Bring the cluster up",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		cfg, _ := config.Yaml()

		compoundAddr.Init(cfg.Cluster.Nodes.Master)

		for i, node := range cfg.Cluster.Nodes.Master {
			dockerClient, _ := docker.Client(address.SSH(cfg.Cluster.Owner, node.Address.External, 0))
			dockerClient.NegotiateAPIVersion(ctx)

			// start zookeeper
			err := container.Up(container.Zookeeper{
				Ctx:          ctx,
				DockerClient: dockerClient,
				Env: env.Zookeeper{
					MYID:    i + 1,
					SERVERS: compoundAddr.Servers,
				},
			})
			if err != nil {
				fmt.Println(err.Error())
			}

			// start mesos master
			err = container.Up(container.MesosMaster{
				Ctx:          ctx,
				DockerClient: dockerClient,
				Env: env.MesosMaster{
					MESOS_HOSTNAME: node.Address.Internal,
					MESOS_IP:       node.Address.Internal,
					MESOS_ZK:       fmt.Sprintf("%s/mesos", compoundAddr.Zookeeper),
				},
			})
			if err != nil {
				fmt.Println(err.Error())
			}

			// start marathon
			err = container.Up(container.Marathon{
				Ctx:          ctx,
				DockerClient: dockerClient,
				Env: env.Marathon{
					MARATHON_HOSTNAME:      node.Address.Internal,
					MARATHON_HTTPS_ADDRESS: node.Address.Internal,
					MARATHON_HTTP_ADDRESS:  node.Address.Internal,
					MARATHON_MASTER:        fmt.Sprintf("%s/mesos", compoundAddr.Zookeeper),
					MARATHON_ZK:            fmt.Sprintf("%s/marathon", compoundAddr.Zookeeper),
				},
			})
			if err != nil {
				fmt.Println(err.Error())
			}

			// start chronos
			err = container.Up(container.Chronos{
				Ctx:          ctx,
				DockerClient: dockerClient,
				Env: env.Chronos{
					CHRONOS_MASTER:   fmt.Sprintf("%s/mesos", compoundAddr.Zookeeper),
					CHRONOS_ZK_HOSTS: compoundAddr.Zookeeper,
				},
			})
			if err != nil {
				fmt.Println(err.Error())
			}

		}

		for _, node := range cfg.Cluster.Nodes.Slave {
			dockerClient, _ := docker.Client(address.SSH(cfg.Cluster.Owner, node.Address.External, 0))
			dockerClient.NegotiateAPIVersion(ctx)

			// start mesos slave
			err := container.Up(container.MesosSlave{
				Ctx:          ctx,
				DockerClient: dockerClient,
				Env: env.MesosSlave{
					MESOS_HOSTNAME: node.Address.Internal,
					MESOS_IP:       node.Address.Internal,
					MESOS_MASTER:   fmt.Sprintf("%s/mesos", compoundAddr.Zookeeper),
				},
			})
			if err != nil {
				fmt.Println(err.Error())
			}

		}

		return nil
	},
}
