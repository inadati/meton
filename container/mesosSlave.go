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

type MesosSlave struct{
	Ctx context.Context
	DockerClient *client.Client
	Env env.MesosSlave
}

func (m MesosSlave) up() error {
	imageName := "mesoscloud/mesos-slave"
	containerName := "mesos-slave"

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
				fmt.Sprintf("MESOS_MASTER=%s", m.Env.MESOS_MASTER),
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

	if err := m.DockerClient.ContainerStart(m.Ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start %v container", containerName)
	}

	// fmt.Println(resp.ID)
	return nil
}
