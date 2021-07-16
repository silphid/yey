package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/silphid/yey/src/cmd"
	yey "github.com/silphid/yey/src/internal"
	"github.com/silphid/yey/src/internal/docker"

	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	options := Options{Remove: new(bool)}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Runs container using given context",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flag("rm").Changed {
				options.Remove = nil
			}
			return run(cmd.Context(), args, options)
		},
	}

	cmd.Flags().BoolVar(options.Remove, "rm", false, "remove container upon exit")
	cmd.Flags().BoolVar(&options.Reset, "reset", false, "remove previous container before starting a fresh one")
	cmd.Flags().BoolVar(&options.Pull, "pull", false, "force pulling image from registry before running")

	return cmd
}

type Options struct {
	Remove *bool
	Reset  bool
	Pull   bool
}

func run(ctx context.Context, names []string, options Options) error {
	contexts, err := yey.LoadContexts()
	if err != nil {
		return err
	}

	lastNames, err := cmd.LoadLastNames()
	if err != nil {
		return err
	}

	names, err = cmd.GetOrPromptContexts(contexts.Context, names, lastNames)
	if err != nil {
		return err
	}

	err = cmd.SaveLastNames(names)
	if err != nil {
		return err
	}

	yeyContext, err := contexts.GetContext(names)
	if err != nil {
		return fmt.Errorf("failed to get context: %w", err)
	}
	if options.Remove != nil {
		yeyContext.Remove = options.Remove
	}

	if yeyContext.Image == "" {
		var err error
		yeyContext.Image, err = readAndBuildDockerfile(ctx, yeyContext.Build, options)
		if err != nil {
			return fmt.Errorf("failed to build yey context image: %w", err)
		}
		yey.Log("using image: %s", yeyContext.Image)
	}

	yey.Log("context:\n--\n%v--", yeyContext)

	// Container name
	containerName := yey.ContainerName(contexts.Path, yeyContext)
	yey.Log("container: %s", containerName)

	// Reset
	if options.Reset {
		yey.Log("removing container first")
		if err := docker.Remove(ctx, containerName, docker.RemoveOptions{}); err != nil {
			return fmt.Errorf("failed to remove container %q: %w", containerName, err)
		}
	}

	// Working directory
	var runOptions docker.RunOptions
	workDir, err := getContainerWorkDir(yeyContext)
	if err != nil {
		return err
	}
	if workDir != "" {
		runOptions.WorkDir = workDir
	}
	yey.Log("working directory: %s", workDir)

	// Pull image first?
	if options.Pull || shouldPull(yeyContext.Image) {
		// TODO: Check running containers with image
		// and prompt for killing it
		yey.Log("Pulling %s", yeyContext.Image)
		docker.Pull(ctx, yeyContext.Image)
	}

	// Banner
	if !yey.IsDryRun {
		if err := ShowBanner(yeyContext.Name); err != nil {
			return err
		}
	}

	return docker.Run(ctx, yeyContext, containerName, runOptions)
}

var tagRegex = regexp.MustCompile(`.*/.*:(.*)`)

func getTagFromImageName(imageName string) string {
	groups := tagRegex.FindStringSubmatch(imageName)
	if groups == nil {
		return ""
	}
	return groups[1]
}

// shouldPull returns whether image should be pulled before running it.
// It returns true when image tag is `latest` or when it is not specified.
func shouldPull(imageName string) bool {
	tag := getTagFromImageName(imageName)
	return tag == "" || tag == "latest"
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

func readAndBuildDockerfile(ctx context.Context, build yey.DockerBuild, options Options) (string, error) {
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
