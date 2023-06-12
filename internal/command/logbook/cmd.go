package logbook

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/sergiught/work-pilot-cli/internal/work"
)

// NewCommand initializes the logbook command.
func NewCommand(repository *work.Repository) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logbook",
		Aliases: []string{"log"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			model := NewModel(repository)

			program := tea.NewProgram(model)
			_, err := program.Run()

			return err
		},
	}

	return cmd
}
