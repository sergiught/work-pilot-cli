package base

import (
	"github.com/spf13/cobra"

	"github.com/sergiught/work-pilot-cli/internal/command/work"
)

type Command struct {
	cmd *cobra.Command
}

func New() *Command {
	cmd := &cobra.Command{
		Use:     "wp",
		Short:   "",
		Long:    "",
		Example: "",
		Version: "dev",
		// SilenceErrors: true,
		// SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(work.New().Command())

	return &Command{
		cmd: cmd,
	}
}

func (c *Command) Command() *cobra.Command {
	return c.cmd
}
