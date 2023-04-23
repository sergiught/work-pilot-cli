package work

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"time"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "work",
		Aliases: []string{"wk"},
		Short:   "",
		Long:    "",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			program := tea.NewProgram(NewWork())
			work, err := program.Run()
			if err != nil {
				return err
			}

			model, ok := work.(Model)
			if !ok {
				return fmt.Errorf("failed to cast model to a work.Model, type is %T", model)
			}

			log.Info("Sleeping...", "seconds", model.Choice)
			time.Sleep(time.Millisecond * time.Duration(model.Choice))
			log.Info("Work done!")

			return nil
		},
	}

	return cmd
}
