package work

import (
	"strconv"

	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/sergiught/work-pilot-cli/internal/work"
)

func NewCommand(repository *work.Repository) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "work",
		Args:    cobra.MaximumNArgs(1),
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			model := NewWorkModel(repository)

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
