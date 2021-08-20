package container

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/meton888/meton/env"
)

type MesosMaster struct{
	Ctx context.Context
	DockerClient *client.Client
	Env env.MesosMaster
}

func (m MesosMaster) up() error {
	imageName := "mesoscloud/mesos-master"
	containerName := "mesos-master"

	out, err := m.DockerClient.ImagePull(m.Ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull %v image", imageName)
	}
	io.Copy(os.Stdout, out)

	resp, err := m.DockerClient.ContainerCreate(
		m.Ctx,
		&container.Config{
			Image: imageName,
			Env: []string{
				fmt.Sprintf("MESOS_HOSTNAME=%s", m.Env.MESOS_HOSTNAME),
				fmt.Sprintf("MESOS_IP=%s", m.Env.MESOS_IP),
				fmt.Sprintf("MESOS_ZK=%s", m.Env.MESOS_ZK),
				"MESOS_PORT=5050",
				"MESOS_LOG_DIR=/var/log/mesos",
				"MESOS_QUORUM=1",
				"MESOS_REGISTRY=in_memory",
				"MESOS_WORK_DIR=/var/lib/mesos",
			},
			// Cmd: []string{"/bin/sh", "-c", "while :; do sleep 10; done"},
		},
		&container.HostConfig{
			NetworkMode: "host",
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
		},
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return fmt.Errorf("failed to create %v container", containerName)
	}

	if err := m.DockerClient.ContainerStart(m.Ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start %v container", containerName)
	}

	// fmt.Println(resp.ID)
	return nil
}
