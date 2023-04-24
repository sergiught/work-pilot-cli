package base

import (
	"fmt"

	"github.com/spf13/cobra"

	logbookCommand "github.com/sergiught/work-pilot-cli/internal/command/logbook"
	workCommand "github.com/sergiught/work-pilot-cli/internal/command/work"
	"github.com/sergiught/work-pilot-cli/internal/platform/database"
	"github.com/sergiught/work-pilot-cli/internal/work"
)

func NewCommand() *cobra.Command {
	workRepository := &work.Repository{}

	cmd := &cobra.Command{
		Use:           "wp",
		Short:         "",
		Long:          "",
		Example:       "",
		Version:       "dev",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.Connect()
			if err != nil {
				return fmt.Errorf("failed to connect to the database: %w", err)
			}

			workRepository.Database = db

			return nil
		},
	}

	cobra.EnableCommandSorting = false
	cmd.AddCommand(
		workCommand.NewCommand(workRepository),
		logbookCommand.NewCommand(workRepository),
	)

	return cmd
}
