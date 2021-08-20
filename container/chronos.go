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

type Chronos struct{
	Ctx context.Context
	DockerClient *client.Client
	Env env.Chronos
}

func (c Chronos) up() error {
	imageName := "mesoscloud/chronos"
	containerName := "chronos"

	out, err := c.DockerClient.ImagePull(c.Ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull %v image", imageName)
	}
	io.Copy(os.Stdout, out)

	resp, err := c.DockerClient.ContainerCreate(
		c.Ctx,
		&container.Config{
			Image: imageName,
			Env: []string{
				"CHRONOS_HTTP_PORT=4400",
				fmt.Sprintf("CHRONOS_MASTER=%s", c.Env.CHRONOS_MASTER),
				fmt.Sprintf("CHRONOS_ZK_HOSTS=%s", c.Env.CHRONOS_ZK_HOSTS),
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

	if err := c.DockerClient.ContainerStart(c.Ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start %v container", containerName)
	}

	return nil
}
