package core

import (
	"github.com/silphid/yey/cli/src/internal/contain"
	"github.com/silphid/yey/cli/src/internal/ctx"
	"github.com/silphid/yey/cli/src/internal/statefile"
)

func (c Core) Start(contextName string) error {
	state, err := statefile.Load(c.homeDir)
	if err != nil {
		return err
	}
	context, err := c.getOrPromptContext(state, contextName)
	if err != nil {
		return err
	}
	return contain.Start(context, state.ImageTag)
}

func (c Core) getOrPromptContext(state statefile.State, name string) (ctx.Context, error) {
	if name == "" {
		name = state.CurrentContext
		if name == "" {
			var err error
			name, err = c.promptContext()
			if err != nil {
				return ctx.None, err
			}
		}
	} else {
		err := c.validateContextName(name)
		if err != nil {
			return ctx.None, err
		}
	}
	return c.GetContext(name)
}