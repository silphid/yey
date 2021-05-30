package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	yey "github.com/silphid/yey/src/internal"
	"github.com/silphid/yey/src/internal/docker"
	"github.com/silphid/yey/src/internal/logging"

	"github.com/silphid/yey/src/cmd"

	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	options := Options{Remove: new(bool)}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs container using given context",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) == 1 {
				name = args[0]
			}
			if !cmd.Flag("rm").Changed {
				options.Remove = nil
			}
			return run(cmd.Context(), name, options)
		},
	}

	cmd.Flags().BoolVar(options.Remove, "rm", false, "remove container upon exit")
	cmd.Flags().BoolVar(&options.Reset, "reset", false, "remove previous container before starting a fresh one")

	return cmd
}

type Options struct {
	Remove *bool
	Reset  bool
}

func run(ctx context.Context, name string, options Options) error {
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	if name == "" {
		var err error
		name, err = cmd.PromptContext(contexts)
		if err != nil {
			return fmt.Errorf("failed to prompt for desired context: %w", err)
		}
	}

	yeyContext, err := contexts.GetContext(name)
	if err != nil {
		return fmt.Errorf("failed to get context with name %q: %w", name, err)
	}
	if options.Remove != nil {
		yeyContext.Remove = options.Remove
	}

	if yeyContext.Image == "" {
		var err error
		yeyContext.Image, err = readAndBuildDockerfile(ctx, yeyContext.Build)
		if err != nil {
			return fmt.Errorf("failed to build yey context image: %w", err)
		}
		logging.Log("using image: %s", yeyContext.Image)
	}

	containerName := yey.ContainerName(contexts.Path, yeyContext)

	if options.Reset {
		if err := docker.Remove(ctx, containerName); err != nil {
			return fmt.Errorf("failed to remove container %q: %w", containerName, err)
		}
	}

	workDir, err := getContainerWorkDir(yeyContext)
	if err != nil {
		return err
	}

	var runOptions []docker.RunOption
	if workDir != "" {
		runOptions = append(runOptions, docker.WithWorkdir(workDir))
	}

	return docker.Start(ctx, yeyContext, containerName, runOptions...)
}

func getContainerWorkDir(yeyContext yey.Context) (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for key, value := range yeyContext.Mounts {
		// Where is work dir relatively to mount dir?
		subDir, err := filepath.Rel(key, workDir)
		if err != nil {
			return "", err
		}

		// Is work dir within mount dir?
		if !strings.HasPrefix(subDir, fmt.Sprintf("..%c", filepath.Separator)) {
			return filepath.Join(value, subDir), nil
		}
	}

	return "", nil
}

func readAndBuildDockerfile(ctx context.Context, build yey.DockerBuild) (string, error) {
	dockerBytes, err := os.ReadFile(build.Dockerfile)
	if err != nil {
		return "", fmt.Errorf("failed to read dockerfile: %w", err)
	}

	imageName := yey.ImageName(dockerBytes)

	if err := docker.Build(ctx, build.Dockerfile, imageName, build.Args, build.Context); err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}

	return imageName, nil
}
