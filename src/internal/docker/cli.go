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

func Start(ctx context.Context, yeyCtx yey.Context, containerName string) error {
	// Determine whether we need to run or exec container
	status, err := getContainerStatus(ctx, containerName)
	if err != nil {
		return err
	}

	switch status {
	case "":
		return runContainer(ctx, yeyCtx, containerName)
	case "exited":
		return startContainer(ctx, containerName)
	case "running":
		return execContainer(ctx, containerName, yeyCtx.Cmd)
	default:
		return fmt.Errorf("container %q in unexpected state %q", containerName, status)
	}
}

func ListContainers(ctx context.Context) ([]string, error) {
	cmd := exec.Command("docker", "ps", "--all", "--filter", "name=yey-*", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(output), "\n"), nil
}

func getContainerStatus(ctx context.Context, name string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "inspect", name, "--format", "{{.State.Status}}")

	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "No such object") {
			return "", nil
		}
		return "", fmt.Errorf("failed to get container status:  %s: %w", output, err)
	}

	return strings.TrimSpace(string(output)), nil
}

func runContainer(ctx context.Context, yeyCtx yey.Context, containerName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	args := []string{
		"run",
		"--name", containerName,
		"-it",
		// "--env LS_COLORS",
		// "--env TERM",
		// "--env TERM_COLOR",
		// "--env TERM_PROGRAM",
		"--env", "YEY_WORK_DIR=" + cwd,
		"--env", "YEY_CONTEXT=" + yeyCtx.Name,
	}

	// Context env vars
	for name, value := range yeyCtx.Env {
		args = append(args, "--env", fmt.Sprintf("%s=%s", name, value))
	}

	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("failed to detect user home directory: %w", err)
	}

	for key, value := range yeyCtx.Mounts {
		args = append(
			args,
			"--volume",
			fmt.Sprintf("%s:%s", strings.ReplaceAll(key, "$HOME", home), value),
		)
	}

	args = append(args, yeyCtx.Image)
	args = append(args, yeyCtx.Cmd...)

	return attachStdPipes(exec.CommandContext(ctx, "docker", args...)).Run()
}

func startContainer(ctx context.Context, containerName string) error {
	return attachStdPipes(exec.CommandContext(ctx, "docker", "start", "-i", containerName)).Run()
}

func execContainer(ctx context.Context, containerName string, cmd []string) error {
	args := append([]string{"exec", "-ti", containerName}, cmd...)
	return attachStdPipes(exec.CommandContext(ctx, "docker", args...)).Run()
}

func attachStdPipes(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}
