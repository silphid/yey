package docker

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
	yey "github.com/silphid/yey/src/internal"
	"github.com/silphid/yey/src/internal/logging"
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

type removeOption struct {
	force bool
}

type RemoveOption func(ro *removeOption)

var ForceRemove RemoveOption = func(ro *removeOption) { ro.force = true }

func Remove(ctx context.Context, containerName string, options ...RemoveOption) error {
	status, err := getContainerStatus(ctx, containerName)
	if err != nil {
		return err
	}

	if status == "" {
		return nil
	}

	return RemoveMany(ctx, []string{containerName}, options...)
}

func RemoveMany(ctx context.Context, containers []string, options ...RemoveOption) error {
	if len(containers) == 0 {
		return nil
	}

	var opts removeOption
	for _, opt := range options {
		opt(&opts)
	}

	args := []string{"rm", "-v"}
	if opts.force {
		args = append(args, "-f")
	}
	args = append(args, containers...)

	return attachStdPipes(exec.CommandContext(ctx, "docker", args...)).Run()
}

func Build(ctx context.Context, dockerPath string, imageTag string, buildArgs map[string]string, context string) error {
	exists, err := imageExists(ctx, imageTag)
	if err != nil {
		return fmt.Errorf("failed to look up image tag %q: %w", imageTag, err)
	}
	if exists {
		logging.Log("found prebuilt image: %q: skipping build step", imageTag)
		return nil
	}

	args := []string{"build", "-f", dockerPath, "-t", imageTag}
	for key, value := range buildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%q", key, value))
	}
	if context == "" {
		context = filepath.Dir(dockerPath)
	}
	args = append(args, context)

	return attachStdPipes(exec.CommandContext(ctx, "docker", args...)).Run()
}

var newlines = regexp.MustCompile(`\r?\n`)

func ListContainers(ctx context.Context) ([]string, error) {
	cmd := exec.Command("docker", "ps", "--all", "--filter", "name=yey-*", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	output = bytes.TrimSpace(output)
	return newlines.Split(string(output), -1), nil
}

func imageExists(ctx context.Context, tag string) (bool, error) {
	output, err := exec.CommandContext(ctx, "docker", "image", "inspect", tag).Output()
	if string(bytes.TrimSpace(output)) == "[]" {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
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
		"-it",
		"--name", containerName,
		"--env", "LS_COLORS",
		"--env", "TERM",
		"--env", "TERM_COLOR",
		"--env", "TERM_PROGRAM",
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

	if yeyCtx.Remove != nil && *yeyCtx.Remove {
		args = append(args, "--rm")
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
