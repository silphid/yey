package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mitchellh/go-homedir"
	yey "github.com/silphid/yey/src/internal"
)

type CLI struct{}

func NewCLI() CLI {
	return CLI{}
}

func (c CLI) Start(ctx context.Context, yeyCtx yey.Context, containerName string) error {
	// Determine whether we need to run or exec container
	isRunning, err := isContainerRunning(containerName)
	if err != nil {
		return err
	}
	args := []string{"run"}
	if !isRunning {
		args = append([]string{"exec"}, args...)
	}

	// Determine extra docker arguments
	extraArgs, err := getDockerArgs(yeyCtx, containerName, !isRunning)
	if err != nil {
		return err
	}
	args = append(args, extraArgs...)

	// Run docker command
	args = append(args, yeyCtx.Cmd...)
	cmd := exec.Command("docker", args...)
	return cmd.Run()
}

func (cli CLI) ListContainers(ctx context.Context) ([]string, error) {
	cmd := exec.Command("docker", "ps", "--all", "--filter", "name=yey-*", "--format", "{{.Name}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return getOutputLines(output), nil
}

func isContainerRunning(name string) (bool, error) {
	cmd := exec.Command("docker", "ps", "--filter", "name="+name, "--format", "foo")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	lines := getOutputLines(output)
	return len(lines) > 0, nil
}

func getOutputLines(value []byte) []string {
	return strings.Split(strings.ReplaceAll(string(value), "\r\n", "\n"), "\n")
}

func getDockerArgs(yeyCtx yey.Context, containerName string, isExec bool) ([]string, error) {
	// Common args
	args := []string{"-it", "--env LS_COLORS", "--env TERM", "--env TERM_COLOR", "--env TERM_PROGRAM"}
	if !isExec {
		args = append(args, "--rm", "--name", containerName)
	}

	// Volumes
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("failed to detect user home directory: %w", err)
	}
	if !isExec {
		for source, target := range yeyCtx.Mounts {
			source = strings.ReplaceAll(source, "$HOME", home)
			args = append(args, "--volume", fmt.Sprintf("%s=%s", source, target))
		}
	}

	// Context env vars
	for name, value := range yeyCtx.Env {
		args = append(args, "--env", fmt.Sprintf("%s=%s", name, value))
	}

	// Built-in env vars
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	args = append(args, "--env", "YEY_WORK_DIR="+cwd)
	args = append(args, "--env", "YEY_CONTEXT="+yeyCtx.Name)

	return args, nil
}
