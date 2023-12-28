package work

import (
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
			selectedTask, err := cmd.Flags().GetString("task")
			if err != nil {
				return err
			}

			selectedTimeDuration, err := cmd.Flags().GetDuration("time")
			if err != nil {
				return err
			}

			model := NewModel(repository)

			if selectedTask != "" && selectedTimeDuration != 0 {
				model.task = selectedTask
				model.timeRemaining = selectedTimeDuration
				model.state = progressView
			}

			program := tea.NewProgram(model)
			_, err = program.Run()

			return err
		},
	}

	cmd.Flags().String("task", "", "The name of the task you want to work on.")
	cmd.Flags().Duration("time", 0, "The amount of time to perform the task for. The value should include a number followed by a unit. Valid units are 's' for seconds, 'm' for minutes, and 'h' for hours. Examples: '3s' for 3 seconds, '5m' for 5 minutes, '1h' for 1 hour.")
	cmd.MarkFlagsRequiredTogether("task", "time")

	return cmd
}
