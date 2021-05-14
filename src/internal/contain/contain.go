package contain

import (
	"fmt"

	"github.com/silphid/yey/cli/src/internal/ctx"
	"github.com/silphid/yey/cli/src/internal/docker/api"
)

const (
	yeyContainerPrefix = "yey"
)

func Start(c ctx.Context, imageTag string) error {
	if c.Container == "" {
		c.Container = "yey"
	}

	containerName := fmt.Sprintf("%s-%s-%s-%s", yeyContainerPrefix, c.Container, c.Name, imageTag)

	docker := api.API{}
	// docker := cli.CLI{}
	return docker.Start(c, imageTag, containerName)
}
