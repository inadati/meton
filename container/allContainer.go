package container

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type AllRecipe struct{}

var All = &AllRecipe{}

func (r *AllRecipe) Down(ctx context.Context, dockerClient *client.Client) error {
	containers, err := dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return fmt.Errorf("failed to get container list")
	}
	for _, container := range containers {
		// fmt.Println("Stopping container ", container.ID[:10], "... ")
		if err := dockerClient.ContainerStop(ctx, container.ID, nil); err != nil {
			return fmt.Errorf("failed to stopped %v container", container.Names[0])
		}
		fmt.Printf("stopped %v container\n", container.Names[0])

		// fmt.Println("Removing container ", container.ID[:10], "... ")
		if err := dockerClient.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			return fmt.Errorf("failed to removed %v container", container.Names[0])
		}

		fmt.Printf("removed %v container\n", container.Names[0])
	}
	return nil
}
