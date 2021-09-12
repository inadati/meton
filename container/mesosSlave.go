package container

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/meton888/meton/env"
)

type MesosSlaveRecipe struct{}

var slave = &MesosSlaveRecipe{}

func (r *MesosSlaveRecipe) Up(ctx context.Context, dockerClient *client.Client, e env.MesosSlave) error {
	imageName := "meton/mesos-slave:1.9.0-centos-7"
	containerName := "mesos-slave"

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
				fmt.Sprintf("MESOS_HOSTNAME=%s", e.MESOS_HOSTNAME),
				fmt.Sprintf("MESOS_IP=%s", e.MESOS_IP),
				fmt.Sprintf("MESOS_MASTER=%s", e.MESOS_MASTER),
				"MESOS_CONTAINERIZERS=docker,mesos",
			},
			// Cmd: []string{"/bin/sh", "-c", "while :; do sleep 10; done"},
		},
		&container.HostConfig{
			NetworkMode: "host",
			PidMode:     "host",
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
			Mounts: []mount.Mount{
				mount.Mount{
					Type:   mount.TypeBind,
					Source: "/usr/bin/docker",
					Target: "/usr/bin/docker",
				},
				mount.Mount{
					Type:   mount.TypeBind,
					Source: "/dev",
					Target: "/dev",
				},
				mount.Mount{
					Type:   mount.TypeBind,
					Source: "/var/run/docker.sock",
					Target: "/var/run/docker.sock",
				},
				mount.Mount{
					Type:   mount.TypeBind,
					Source: "/var/log/mesos",
					Target: "/var/log/mesos",
				},
				mount.Mount{
					Type:   mount.TypeBind,
					Source: "/tmp/mesos",
					Target: "/tmp/mesos",
				},
			},
			Privileged: true,
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
