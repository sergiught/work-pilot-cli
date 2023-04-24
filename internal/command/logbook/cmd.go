package logbook

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sergiught/work-pilot-cli/internal/command/work"
	"github.com/spf13/cobra"
)

func NewCommand(repository *work.Repository) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logbook",
		Aliases: []string{"log"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			var workItems []work.Work

			repository.Database.Find(&workItems)
			if repository.Database.Error != nil {
				return repository.Database.Error
			}

			model := NewModel(workItems)
			program := tea.NewProgram(model)
			_, err := program.Run()
			if err != nil {
				return err
			}

			return err
		},
	}

	cmd.Flags().String("task", "", "The name of the task you want to work on")

	return cmd
}
