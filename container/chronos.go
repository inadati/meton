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

type ChronosRecipe struct {}

var Chronos = &ChronosRecipe{}

func (r *ChronosRecipe) Up(ctx context.Context, dockerClient *client.Client, e env.Chronos) error {
	imageName := "mesoscloud/chronos"
	containerName := "chronos"

	out, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull %v image", imageName)
	}
	io.Copy(os.Stdout, out)

	resp, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
			Env: []string{
				"CHRONOS_HTTP_PORT=4400",
				fmt.Sprintf("CHRONOS_MASTER=%s", e.CHRONOS_MASTER),
				fmt.Sprintf("CHRONOS_ZK_HOSTS=%s", e.CHRONOS_ZK_HOSTS),
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

	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start %v container", containerName)
	}

	return nil
}
