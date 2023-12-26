package work

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/sergiught/work-pilot-cli/internal/work"
)

// NewCommand initializes the work command.
func NewCommand(repository *work.Repository) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "work",
		Args:  cobra.NoArgs,
		Short: "Track time spent on a work task",
		Long:  "This command initializes time tracking for the specified work task.",
		Example: `  wp work
  wp work --task "cooking"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			model := NewModel(repository)

			task, err := cmd.Flags().GetString("task")
			if err != nil {
				return err
			}

			model.task = task

			if len(args) > 0 {
				choice, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}

				model.choice = choice
				model.timeRemaining = choice
			}

			program := tea.NewProgram(model)
			_, err = program.Run()

			return err
		},
	}

	cmd.Flags().String("task", "", "The name of the task you want to work on")

	return cmd
}
