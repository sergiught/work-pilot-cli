package work

import (
	"time"

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

			selectedTask, err := cmd.Flags().GetString("task")
			if err != nil {
				return err
			}

			selectedTimeDuration, err := cmd.Flags().GetDuration("time")
			if err != nil {
				return err
			}

			if selectedTask != "" {
				model.task = selectedTask
			}

			if selectedTimeDuration != 0 {
				model.timeRemaining = selectedTimeDuration
			}

			program := tea.NewProgram(model)
			_, err = program.Run()

			return err
		},
	}

	cmd.Flags().String("task", "", "The name of the task you want to work on")
	cmd.Flags().Duration("time", time.Minute*0, "The amount of time to perform the task for")

	return cmd
}
