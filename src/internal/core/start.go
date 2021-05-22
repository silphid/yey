package core

import (
	yey "github.com/silphid/yey/src/internal"
	"github.com/silphid/yey/src/internal/contain"
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
		var err error
		name, err = c.promptContext()
		if err != nil {
			return yey.Context{}, err
		}
	} else {
		err := c.validateContextName(name)
		if err != nil {
			return yey.Context{}, err
		}
	}
	return c.GetContext(name)
}
