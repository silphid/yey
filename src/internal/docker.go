package yey

type Docker interface {
	Start(c Context, imageTag, containerName string) error
}
