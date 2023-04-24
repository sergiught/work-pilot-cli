package logbook

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sergiught/work-pilot-cli/internal/command/work"
	"strconv"
	"time"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	table table.Model
}

func NewModel(workItems []work.Work) *Model {
	columns := []table.Column{
		{Title: "Work Task", Width: 30},
		{Title: "Duration", Width: 10},
		{Title: "Date", Width: 11},
		{Title: "Started at", Width: 11},
		{Title: "Finished at", Width: 11},
	}

	var rows []table.Row
	for _, item := range workItems {
		rows = append(rows, table.Row{
			item.Task,
			strconv.Itoa(item.Duration),
			item.CreatedAt.Format("2006-01-02"),
			item.CreatedAt.Format("15:04:05"),
			item.CreatedAt.Add(time.Duration(item.Duration) * time.Second).Format("15:04:05"),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)

	return &Model{
		table: t,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m Model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}
