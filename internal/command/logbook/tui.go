package logbook

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"github.com/sergiught/work-pilot-cli/internal/work"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	repository *work.Repository

	table table.Model
}

func NewModel(repository *work.Repository) *Model {
	workTasks, err := repository.GetAllWorkTasks()
	if err != nil {
		log.Error("failed to get all work tasks from the database", err)
		return nil
	}

	var rows []table.Row
	for _, task := range workTasks {
		rows = append(rows, table.Row{
			task.Name,
			strconv.Itoa(task.Duration),
			task.CreatedAt.Format("2006-01-02"),
			task.CreatedAt.Format("15:04:05"),
			task.CreatedAt.Add(time.Duration(task.Duration) * time.Second).Format("15:04:05"),
		})
	}

	columns := []table.Column{
		{Title: "Work Task", Width: 30},
		{Title: "Duration", Width: 10},
		{Title: "Date", Width: 11},
		{Title: "Started at", Width: 11},
		{Title: "Finished at", Width: 11},
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
		repository: repository,
		table:      t,
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
