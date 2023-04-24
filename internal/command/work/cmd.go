package work

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "work",
		Aliases: []string{"wk"},
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			program := tea.NewProgram(NewWorkModel())
			work, err := program.Run()
			if err != nil {
				return err
			}

			model, ok := work.(Model)
			if !ok {
				return fmt.Errorf("failed to cast model to a work.Model, type is %T", model)
			}

			return nil
		},
	}

	return cmd
}
