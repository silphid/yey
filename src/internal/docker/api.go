package docker

import (
	"context"
	"fmt"
	"strings"

	yey "github.com/silphid/yey/src/internal"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/silphid/yey/src/internal/logging"

	homedir "github.com/mitchellh/go-homedir"
)

func NewAPI() (API, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return API{}, err
	}
	return API{client}, nil
}

type API struct {
	client *client.Client
}

func (api API) Start(ctx context.Context, yeyCtx yey.Context, containerName string) error {
	container, err := api.getContainer(ctx, containerName)
	if err != nil {
		return err
	}

	if container == nil {
		logging.Log("Container %q does not already exist", containerName)
		container, err = api.createContainer(ctx, yeyCtx, containerName)
		if err != nil {
			return err
		}
		logging.Log("Container %q created", containerName)
	} else {
		logging.Log("Reusing existing container %q (%s)", containerName, container.State)
	}

	if container.State != "running" {
		logging.Log("Starting container %q", containerName)
		err = api.client.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
		if err != nil {
			return err
		}
	} else {
		logging.Log("Executing in container %q", containerName)
		// TODO
	}

	return nil
}

func (api API) getContainer(ctx context.Context, name string) (*types.Container, error) {
	containers, err := api.client.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", name)),
	})
	if err != nil {
		return nil, err
	}

	if len(containers) > 0 {
		return &containers[0], nil
	}

	return nil, nil
}

func (api API) createContainer(ctx context.Context, c yey.Context, name string) (*types.Container, error) {
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
	_, err = api.client.ContainerCreate(ctx, &config, &hostConfig, &networkingConfig, nil, name)
	if err != nil {
		return nil, err
	}

	return api.getContainer(ctx, name)
}

// TODO should return strings since library code should not be expected to handle the printing
func (api API) ListContainers(ctx context.Context) error {
	containers, err := api.client.ContainerList(ctx, types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("name", "yey-*")),
	})
	if err != nil {
		return err
	}

	for _, container := range containers {
		// HMMMMMMMMM SUSPICIOUS
		fmt.Println(strings.TrimPrefix(container.Names[0], "/"))
	}
	return nil
}
