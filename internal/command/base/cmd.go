package base

import (
	"github.com/spf13/cobra"

	"github.com/sergiught/work-pilot-cli/internal/command/work"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "wp",
		Short:         "",
		Long:          "",
		Example:       "",
		Version:       "dev",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(work.NewCommand())

	return cmd
}
