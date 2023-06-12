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

// Model is the TUI model for the logbook command.
type Model struct {
	repository *work.Repository

	table table.Model
}

// NewModel initializes the TUI model for the logbook command.
func NewModel(repository *work.Repository) *Model {
	workTasks, err := repository.GetAllWorkTasks()
	if err != nil {
		log.Error("failed to get all work tasks from the database", err)
		return nil
	}

	var rows []table.Row
	var workTaskColumnLength, durationColumnLength, dateColumnLength, startedAtColumnLength, finishedAtColumnLength int
	for _, task := range workTasks {
		taskName := task.Name
		duration := strconv.Itoa(task.Duration)
		date := task.CreatedAt.Format("2006-01-02")
		startedAt := task.CreatedAt.Format("15:04:05")
		finishedAt := task.CreatedAt.Add(time.Duration(task.Duration) * time.Second).Format("15:04:05")

		if workTaskColumnLength < len(task.Name) {
			workTaskColumnLength = len(task.Name)
		}
		if durationColumnLength < len(strconv.Itoa(task.Duration)) {
			durationColumnLength = len(strconv.Itoa(task.Duration))
		}
		if dateColumnLength < len(task.CreatedAt.Format("2006-01-02")) {
			dateColumnLength = len(task.CreatedAt.Format("2006-01-02"))
		}
		if startedAtColumnLength < len(task.CreatedAt.Format("15:04:05")) {
			startedAtColumnLength = len(task.CreatedAt.Format("15:04:05"))
		}
		if finishedAtColumnLength < len(task.CreatedAt.Add(time.Duration(task.Duration)*time.Second).Format("15:04:05")) {
			finishedAtColumnLength = len(task.CreatedAt.Add(time.Duration(task.Duration) * time.Second).Format("15:04:05"))
		}

		rows = append(rows, table.Row{
			taskName,
			duration,
			date,
			startedAt,
			finishedAt,
		})
	}

	if workTaskColumnLength < len("Work Task") {
		workTaskColumnLength = len("Work Task")
	}
	if durationColumnLength < len("Duration") {
		durationColumnLength = len("Duration")
	}
	if dateColumnLength < len("Date") {
		dateColumnLength = len("Date")
	}
	if startedAtColumnLength < len("Started at") {
		startedAtColumnLength = len("Started at")
	}
	if finishedAtColumnLength < len("Finished at") {
		finishedAtColumnLength = len("Finished at")
	}

	columns := []table.Column{
		{Title: "Work Task", Width: workTaskColumnLength},
		{Title: "Duration", Width: durationColumnLength},
		{Title: "Date", Width: dateColumnLength},
		{Title: "Started at", Width: startedAtColumnLength},
		{Title: "Finished at", Width: finishedAtColumnLength},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(rows)),
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

// Init currently does nothing.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update currently does nothing.
func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

// View holds the view logic for the main logbook Model.
func (m Model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}
