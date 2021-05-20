package contain

import (
	"fmt"
	"regexp"

	"github.com/silphid/yey/src/internal/docker/api"
	"github.com/silphid/yey/src/internal/yey"
)

const (
	yeyContainerPrefix = "yey"
)

func Start(c yey.Context) error {
	shortImageName, err := getShortImageName(c.Image)
	if err != nil {
		return err
	}
	containerName := fmt.Sprintf("%s-%s-%s", yeyContainerPrefix, shortImageName, c.Name)

	docker := api.API{}
	// docker := cli.CLI{}
	return docker.Start(c, containerName)
}

var imageNameRegex = regexp.MustCompile(`(.*/)?(.+?)(:.*)?$`)

// getShortImageName returns short image name without any registry
// prefix or tag suffix, to be used as part of container name.
func getShortImageName(imageName string) (string, error) {
	submatches := imageNameRegex.FindStringSubmatch(imageName)
	if len(submatches) < 4 {
		return "", fmt.Errorf("malformed image name %q", imageName)
	}
	return submatches[2], nil
}
