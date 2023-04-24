package work

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"strconv"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "work",
		Aliases: []string{"wk"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			model := NewWorkModel()

			if len(args) > 0 {
				choice, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}

				model.choice = choice
				model.timeRemaining = choice
			}

			program := tea.NewProgram(model)
			_, err := program.Run()

			return err
		},
	}

	return cmd
}
