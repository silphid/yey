package core

import (
	yey "github.com/silphid/yey/src/internal"
)

type Core struct {
	contexts yey.Contexts
}

func New() (Core, error) {
	contexts, err := yey.ReadAndParseContextFile()
	if err != nil {
		return Core{}, err
	}
	return Core{contexts}, nil
}

// GetContextNames returns the list of all context names user can
// choose from, including the special "base" context.
func (c Core) GetContextNames() ([]string, error) {
	return c.contexts.GetNames(), nil
}

// GetContext finds shared/user base/named contexts and returns their merged result.
// If name is "base", only the merged base context is returned.
func (c Core) GetContext(name string) (yey.Context, error) {
	return c.contexts.GetContext(name)
}
