package base

import (
	"fmt"
	"github.com/sergiught/work-pilot-cli/internal/command/work"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
			database, err := gorm.Open(sqlite.Open("work-pilot.db"), &gorm.Config{})
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}
			err = database.AutoMigrate(&work.Work{})
			if err != nil {
				return fmt.Errorf("failed to migrate database: %w", err)
			}

			workRepository.Database = database

			return nil
		},
	}

	cmd.AddCommand(work.NewCommand(workRepository))

	return cmd
}
