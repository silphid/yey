package docker

import "github.com/silphid/yey/cli/src/internal/ctx"

type Docker interface {
	Start(c ctx.Context, imageTag, containerName string) error
}
