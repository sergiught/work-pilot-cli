package work

import (
	"github.com/spf13/cobra"
)

type Command struct {
	cmd *cobra.Command
}

func New() *Command {
	cmd := &cobra.Command{
		Use:     "work",
		Aliases: []string{"wk"},
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	return &Command{
		cmd: cmd,
	}
}

func (c *Command) Command() *cobra.Command {
	return c.cmd
}
