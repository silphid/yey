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

	yey "github.com/silphid/yey/src/internal"
)

type RunOptions struct {
	WorkDir string
}

func Run(ctx context.Context, yeyCtx yey.Context, containerName string, options RunOptions) error {
	// Determine whether we need to run or exec container
	status, err := getContainerStatus(ctx, containerName)
	if err != nil {
		return err
	}

	switch status {
	case "":
		yey.Log("running new container %q", containerName)
		return runContainer(ctx, yeyCtx, containerName, options)
	case "exited":
		yey.Log("restarting stopped container %q", containerName)
		return startContainer(ctx, containerName, options)
	case "running":
		yey.Log("executing new shell in running container %q", containerName)
		return execContainer(ctx, containerName, yeyCtx.Cmd, options)
	default:
		return fmt.Errorf("container %q in unexpected state %q", containerName, status)
	}
}

type RemoveOptions struct {
	Force bool
}

func Remove(ctx context.Context, containerName string, options RemoveOptions) error {
	status, err := getContainerStatus(ctx, containerName)
	if err != nil {
		return err
	}

	if status == "" {
		return nil
	}

	return RemoveMany(ctx, []string{containerName}, options)
}

func RemoveMany(ctx context.Context, containers []string, options RemoveOptions) error {
	if len(containers) == 0 {
		return nil
	}

	args := []string{"rm", "-v"}
	if options.Force {
		args = append(args, "-f")
	}
	args = append(args, containers...)

	return run(ctx, args...)
}

func Build(ctx context.Context, dockerPath string, imageTag string, buildArgs map[string]string, context string) error {
	exists, err := imageExists(ctx, imageTag)
	if err != nil {
		return fmt.Errorf("failed to look up image tag %q: %w", imageTag, err)
	}
	if exists {
		yey.Log("found prebuilt image: %q: skipping build step", imageTag)
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

	return run(ctx, args...)
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

func runContainer(ctx context.Context, yeyCtx yey.Context, containerName string, options RunOptions) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	args := []string{
		"run",
		"-it",
		"--name", containerName,
		"--env", "YEY_WORK_DIR=" + cwd,
		"--env", "YEY_CONTEXT=" + yeyCtx.Name,
	}

	// Context env vars
	for name, value := range yeyCtx.Env {
		arg := name
		if value != "" {
			arg = fmt.Sprintf("%s=%s", name, value)
		}
		args = append(args, "--env", arg)
	}

	// Mount binds
	for key, value := range yeyCtx.Mounts {
		args = append(
			args,
			"--volume",
			fmt.Sprintf("%s:%s", key, value),
		)
	}

	// Remove container upon exit?
	if yeyCtx.Remove != nil && *yeyCtx.Remove {
		args = append(args, "--rm")
	}

	// Network mode
	network := yeyCtx.Network
	if network == "" {
		network = "host"
	}
	args = append(args, "--network", network)

	// Work directory
	if options.WorkDir != "" {
		args = append(args, "--workdir", options.WorkDir)
	}

	args = append(args, yeyCtx.Image)
	args = append(args, yeyCtx.Cmd...)

	return run(ctx, args...)
}

func startContainer(ctx context.Context, containerName string, options RunOptions) error {
	return run(ctx, "start", "-i", containerName)
}

func execContainer(ctx context.Context, containerName string, cmd []string, options RunOptions) error {
	args := []string{"exec", "-ti"}
	if options.WorkDir != "" {
		args = append(args, "--workdir", options.WorkDir)
	}
	args = append(args, containerName)
	args = append(args, cmd...)

	return run(ctx, args...)
}

func run(ctx context.Context, args ...string) error {
	if yey.IsDryRun {
		fmt.Printf("docker %s\n", strings.Join(args, " "))
		return nil
	}
	return attachStdPipes(exec.CommandContext(ctx, "docker", args...)).Run()
}

func attachStdPipes(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}
