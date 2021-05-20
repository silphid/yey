package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/silphid/yey/cli/src/internal/logging"
	"github.com/silphid/yey/cli/src/internal/yey"
)

type API struct{}

func (a API) Start(c yey.Context, containerName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	container, err := getContainer(cli, containerName)
	if err != nil {
		return err
	}

	if container == nil {
		logging.Log("Container %q does not already exist", containerName)
		container, err = createContainer(cli, c, containerName)
		if err != nil {
			return err
		}
		logging.Log("Container %q created", containerName)
	} else {
		logging.Log("Reusing existing container %q (%s)", containerName, container.State)
	}

	if container.State != "running" {
		logging.Log("Starting container %q", containerName)
		err = cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
		if err != nil {
			return err
		}
	} else {
		logging.Log("Executing in container %q", containerName)
		// TODO
	}

	return nil
}

func getContainer(cli *client.Client, name string) (*types.Container, error) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("name", name)),
	})
	if err != nil {
		return nil, err
	}

	if len(containers) > 0 {
		return &containers[0], nil
	}
	return nil, nil
}

func createContainer(cli *client.Client, c yey.Context, name string) (*types.Container, error) {
	if c.Image == "" {
		return nil, fmt.Errorf("missing required property %q", "image")
	}

	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("failed to detect user home directory: %w", err)
	}

	mounts := make([]mount.Mount, 0)
	for source, target := range c.Mounts {
		source = strings.ReplaceAll(source, "$HOME", home)
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: source,
			Target: target,
		})
	}

	config := container.Config{
		Image:        c.Image,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Entrypoint:   []string{"sh"},
	}
	hostConfig := container.HostConfig{
		Mounts: mounts,
	}
	networkingConfig := network.NetworkingConfig{}

	logging.Log("Creating container from image: %s", c.Image)
	_, err = cli.ContainerCreate(context.Background(), &config, &hostConfig, &networkingConfig, nil, name)
	if err != nil {
		return nil, err
	}

	return getContainer(cli, name)
}

func ListContainers() error {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("name", "yey-*")),
	})
	if err != nil {
		return err
	}

	for _, container := range containers {
		fmt.Println(strings.TrimPrefix(container.Names[0], "/"))
	}
	return nil
}
