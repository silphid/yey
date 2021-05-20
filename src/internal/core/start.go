package core

import (
	"github.com/silphid/yey/cli/src/internal/contain"
	"github.com/silphid/yey/cli/src/internal/ctx"
	"github.com/silphid/yey/cli/src/internal/yey"
)

func (c Core) Start(contextName string) error {
	context, err := c.getOrPromptContext(contextName)
	if err != nil {
		return err
	}
	return contain.Start(context)
}

func (c Core) getOrPromptContext(name string) (yey.Context, error) {
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
