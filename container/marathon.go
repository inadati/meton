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

type MarathonRecipe struct{}

var Marathon = &MarathonRecipe{}

func (r *MarathonRecipe) Up(ctx context.Context, dockerClient *client.Client, e env.Marathon) error {
	imageName := "mesoscloud/marathon"
	containerName := "marathon"

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
				fmt.Sprintf("MARATHON_HOSTNAME=%s", e.MARATHON_HOSTNAME),
				fmt.Sprintf("MARATHON_HTTPS_ADDRESS=%s", e.MARATHON_HTTPS_ADDRESS),
				fmt.Sprintf("MARATHON_HTTP_ADDRESS=%s", e.MARATHON_HTTP_ADDRESS),
				fmt.Sprintf("MARATHON_MASTER=%s", e.MARATHON_MASTER),
				fmt.Sprintf("MARATHON_ZK=%s", e.MARATHON_ZK),
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
