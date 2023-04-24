package work

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"strconv"
)

func NewCommand(repository *Repository) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "work",
		Aliases: []string{"wk"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			task, err := cmd.Flags().GetString("task")
			if err != nil {
				return err
			}

			model := NewWorkModel()
			if len(args) > 0 {
				choice, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}

				model.choice = choice
				model.timeRemaining = choice
			}

			model.task = task

			program := tea.NewProgram(model)
			_, err = program.Run()
			if err != nil {
				return err
			}

			work := Work{
				Task:     model.task,
				Duration: model.choice,
			}

			if model.task == "" {
				work.Task = "generic"
			}

			repository.Database.Create(&work)
			if repository.Database.Error != nil {
				return repository.Database.Error
			}

			return nil
		},
	}

	cmd.Flags().String("task", "", "The name of the task you want to work on")

	return cmd
}
