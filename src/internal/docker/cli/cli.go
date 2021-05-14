package cli

import (
	"fmt"

	"github.com/silphid/yey/cli/src/internal/ctx"
)

type CLI struct{}

func (c CLI) Start(ct ctx.Context, imageTag, containerName string) error {
	return fmt.Errorf("not implemented")

	// Get running ID and state
	// docker ps --all --filter "name=al" --format '{{.ID}}|{{.State}}'
	// state="exited"|"running"

	// Run
	// docker run -it --name al alpine

	// Exec
	// run docker exec -it "${DOCKER_EXEC_ARGS[@]}" "${DOCKER_CONTAINER}" zsh
}

type containerState string

const (
	stateRunning containerState = "running"
	stateExited                 = "exited"
)

type container struct {
	state containerState
}

func getContainer(name string) (container, error) {
	return container{}, nil
}
