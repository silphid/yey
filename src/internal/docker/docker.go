package docker

import "github.com/silphid/yey/src/internal/yey"

type Docker interface {
	Start(c yey.Context, imageTag, containerName string) error
}
