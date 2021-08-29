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

// type Zookeeper struct {
// 	Ctx context.Context
// 	DockerClient *client.Client
// 	Env env.Zookeeper

// }

type ZookeeperRecipe struct{}

var Zookeeper = &ZookeeperRecipe{}

func (r *ZookeeperRecipe) Up(ctx context.Context, dockerClient *client.Client, e env.Zookeeper) error {
	imageName := "mesoscloud/zookeeper"
	containerName := "zookeeper"

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
				fmt.Sprintf("MYID=%d", e.MYID),
				fmt.Sprintf("SERVERS=%s", e.SERVERS),
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

	// fmt.Println(resp.ID)
	return nil
}
