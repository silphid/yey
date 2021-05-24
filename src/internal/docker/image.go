package docker

import (
	"fmt"
	"regexp"
)

var imageNameRegex = regexp.MustCompile(`(.*/)?(.+?)(:.*)?$`)

// GetShortImageName returns short image name without any registry
// prefix or tag suffix, to be used as part of container name.
func GetShortImageName(imageName string) (string, error) {
	submatches := imageNameRegex.FindStringSubmatch(imageName)
	if len(submatches) < 4 {
		return "", fmt.Errorf("malformed image name %q", imageName)
	}
	return submatches[2], nil
}
