package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/silphid/yey/cli/src/internal/cfg"
	"github.com/silphid/yey/cli/src/internal/ctx"
	"github.com/silphid/yey/cli/src/internal/logging"
	"github.com/silphid/yey/cli/src/internal/statefile"
)

const (
	homeVar       = "YEY_HOME"
	homeDirName   = ".yey"
	configDirName = "config"
)

type Core struct {
	homeDir   string
	sharedDir string
}

func New() (Core, error) {
	homeDir, err := getYeyHomeDir()
	if err != nil {
		return Core{}, err
	}
	cloneDir, err := getYeyCloneDir(homeDir)
	if err != nil {
		return Core{}, err
	}
	sharedDir := filepath.Join(cloneDir, configDirName)
	return Core{
		homeDir:   homeDir,
		sharedDir: sharedDir,
	}, nil
}

// GetContextNames returns the list of all context names user can choose from including
// the special "base" and "none" contexts.
func (c Core) GetContextNames() ([]string, error) {
	return cfg.GetContextNames(c.sharedDir, c.homeDir)
}

// GetContext finds shared/user base/named contexts and returns their merged result.
// If name is "base", only the merged base context is returned.
// If name is "none", an empty context is returned.
func (c Core) GetContext(name string) (ctx.Context, error) {
	return cfg.GetContext(c.sharedDir, c.homeDir, name)
}

// GetState loads and returns current state
func (c Core) GetState() (statefile.State, error) {
	return statefile.Load(c.homeDir)
}

// getYeyCloneDir returns the path to the yey clone directory, specified by required YEY_CLONE env var
func getYeyCloneDir(homeDir string) (string, error) {
	state, err := statefile.Load(homeDir)
	if err != nil {
		return "", err
	}
	if state.CloneDir == "" {
		return "", fmt.Errorf("clone dir not defined in state.yaml, please reinstall")
	}
	return state.CloneDir, nil
}

// getYeyHomeDir returns the path to the yey home directory, optionally specified by YEY_HOME env var
// (defaults to ~/.yey)
func getYeyHomeDir() (homeDir string, err error) {
	defer func() {
		if err == nil {
			logging.Log("Using yey home dir: %s", homeDir)
		}
	}()

	homeDir, ok := os.LookupEnv(homeVar)
	if ok && homeDir != "" {
		return
	}

	home, err := homedir.Dir()
	if err != nil {
		err = fmt.Errorf("failed to detect user home directory: %w", err)
		return
	}
	homeDir = filepath.Join(home, homeDirName)
	return
}
