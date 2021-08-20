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

type Marathon struct{
	Ctx context.Context
	DockerClient *client.Client
	Env env.Marathon
}

func (m Marathon) up() error {
	imageName := "mesoscloud/marathon"
	containerName := "marathon"

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
				fmt.Sprintf("MARATHON_HOSTNAME=%s", m.Env.MARATHON_HOSTNAME),
				fmt.Sprintf("MARATHON_HTTPS_ADDRESS=%s", m.Env.MARATHON_HTTPS_ADDRESS),
				fmt.Sprintf("MARATHON_HTTP_ADDRESS=%s", m.Env.MARATHON_HTTP_ADDRESS),
				fmt.Sprintf("MARATHON_MASTER=%s", m.Env.MARATHON_MASTER),
				fmt.Sprintf("MARATHON_ZK=%s", m.Env.MARATHON_ZK),
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

	return nil
}
