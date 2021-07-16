package docker

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
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
		return execContainer(ctx, yeyCtx, containerName, options)
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

func Pull(ctx context.Context, image string) error {
	return run(ctx, "pull", image)
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

func ListContainers(ctx context.Context, all bool) ([]string, error) {
	// Compute args
	args := []string{"ps", "--filter", "name=yey-*", "--format", "{{.Names}}"}
	if all {
		args = append(args, "--all")
	}

	cmd := exec.Command("docker", args...)

	// Parse output
	outputBuf, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: docker %s", strings.Join(args, " "))
	}
	output := string(bytes.TrimSpace(outputBuf))
	if output == "" {
		return []string{}, nil
	}
	containers := newlines.Split(string(output), -1)

	// Sort
	sort.Strings(containers)
	return containers, nil
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
		args = append(args, "--env", fmt.Sprintf("%s=%s", name, value))
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

	// EntryPoint
	if yeyCtx.EntryPoint != "" {
		args = append(args, "--entrypoint", yeyCtx.EntryPoint)
	}

	args = append(args, yeyCtx.Image)
	args = append(args, yeyCtx.Cmd...)

	return run(ctx, args...)
}

func startContainer(ctx context.Context, containerName string, options RunOptions) error {
	return run(ctx, "start", "-i", containerName)
}

func execContainer(ctx context.Context, yeyCtx yey.Context, containerName string, options RunOptions) error {
	args := []string{"exec", "-ti"}
	if options.WorkDir != "" {
		args = append(args, "--workdir", options.WorkDir)
	}
	args = append(args, containerName)

	// Entrypoint/command
	if yeyCtx.EntryPoint == "" && len(yeyCtx.Cmd) == 0 {
		// Default to sh because docker exec requires some command
		args = append(args, "sh")
	} else {
		if yeyCtx.EntryPoint != "" {
			args = append(args, yeyCtx.EntryPoint)
		}
		args = append(args, yeyCtx.Cmd...)
	}

	return run(ctx, args...)
}

func run(ctx context.Context, args ...string) error {
	if yey.IsDryRun || yey.IsVerbose {
		cmd := fmt.Sprintf("docker %s", strings.Join(quoteArgsWithSpecialChars(args), " "))
		if yey.IsDryRun {
			fmt.Println(cmd)
			return nil
		}
		yey.Log(cmd)
	}

	return attachStdPipes(exec.CommandContext(ctx, "docker", args...)).Run()
}

var specialCharsRegex = regexp.MustCompile(`\s`)

func quoteArgsWithSpecialChars(args []string) []string {
	escaped := make([]string, 0, len(args))
	for _, a := range args {
		if specialCharsRegex.MatchString(a) {
			a = `"` + a + `"`
		}
		escaped = append(escaped, a)
	}
	return escaped
}

func attachStdPipes(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}
