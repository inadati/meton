package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type All struct {
	Ctx context.Context
	DockerClient *client.Client
}

func (a All) down() error {
	containers, err := a.DockerClient.ContainerList(a.Ctx, types.ContainerListOptions{})
	if err != nil {
		return fmt.Errorf("failed to get container list")
	}
	for _, container := range containers {
		// fmt.Println("Stopping container ", container.ID[:10], "... ")
		if err := a.DockerClient.ContainerStop(a.Ctx, container.ID, nil); err != nil {
			return fmt.Errorf("failed to stopped %v container", container.Names[0])
		}
		fmt.Printf("stopped %v container\n", container.Names[0])

		// fmt.Println("Removing container ", container.ID[:10], "... ")
		if err := a.DockerClient.ContainerRemove(a.Ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			return fmt.Errorf("failed to removed %v container", container.Names[0])
		}

		fmt.Printf("removed %v container\n", container.Names[0])
	}
	return nil
}
