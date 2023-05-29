package work

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"github.com/sergiught/work-pilot-cli/internal/work"
)

const (
	progressMaxWidth = 80
	progressPadding  = 2
	listWidth        = 20
	listHeight       = 14
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1).ColorWhitespace(true)
	infoTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	progressHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
)

type state int

const (
	taskSelection state = iota
	timeSelection
	customTimeSelection
	progressView
)

type Model struct {
	repository *work.Repository

	taskInput textinput.Model
	timeList  list.Model
	timeInput textinput.Model
	progress  progress.Model

	state      state
	isQuitting bool

	choice        int
	timeRemaining int
	task          string
}

func NewModel(repository *work.Repository) *Model {
	customTaskInput := NewTaskInput()
	timeList := NewTimeList()
	customTimeInput := NewTimeInput()
	progressIndicator := NewProgressIndicator()

	return &Model{
		repository: repository,
		taskInput:  customTaskInput,
		timeList:   timeList,
		timeInput:  customTimeInput,
		progress:   progressIndicator,
		state:      taskSelection,
	}
}

func NewTaskInput() textinput.Model {
	taskInput := textinput.New()
	taskInput.SetValue("Work")
	taskInput.Placeholder = "Work"
	taskInput.Focus()
	taskInput.CharLimit = 256
	taskInput.Width = 256

	return taskInput
}

func NewTimeList() list.Model {
	items := []list.Item{
		listItem{
			label: "20 seconds",
			value: 20,
		},
		listItem{
			label: "40 seconds",
			value: 40,
		},
		listItem{
			label: "60 seconds",
			value: 60,
		},
		listItem{
			label: "Custom Value",
		},
	}

	timeList := list.New(items, itemDelegate{}, listWidth, listHeight)
	timeList.Title = "How many minutes do you want to work for?"
	timeList.SetShowStatusBar(false)
	timeList.SetFilteringEnabled(false)
	timeList.Styles.Title = titleStyle
	timeList.Styles.PaginationStyle = paginationStyle
	timeList.Styles.HelpStyle = helpStyle

	return timeList
}

type listItem struct {
	label string
	value int
}

func (i listItem) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) Update(tea.Msg, *list.Model) tea.Cmd {
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(listItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("• %s", i.label)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func NewTimeInput() textinput.Model {
	timeInput := textinput.New()
	timeInput.Placeholder = "0"
	timeInput.Focus()
	timeInput.CharLimit = 32
	timeInput.Width = 20

	return timeInput
}

func NewProgressIndicator() progress.Model {
	progressIndicator := progress.New(progress.WithDefaultGradient())
	return progressIndicator
}

func (m Model) Init() tea.Cmd {
	if m.choice != 0 {
		return tick()
	}

	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var commands []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.timeList.SetWidth(msg.Width)

		m.progress.Width = msg.Width - progressPadding*2 - 4
		if m.progress.Width > progressMaxWidth {
			m.progress.Width = progressMaxWidth
		}

		return m, nil
	case selectedWorkTask:
		m.task = msg.task
		m.taskInput.Reset()

		m.state = timeSelection
		m.timeList, cmd = m.timeList.Update(msg)

		return m, cmd
	case selectedWorkTimeFromInput:
		m.state = customTimeSelection
		m.timeInput, cmd = m.timeInput.Update(msg)

		return m, cmd
	case selectedWorkTimeFromList:
		m.choice = msg.time
		m.timeRemaining = msg.time
		m.state = progressView

		return m, tick()
	case selectedCustomTime:
		m.choice = msg.time
		m.timeRemaining = msg.time
		m.state = progressView

		return m, tick()
	case workFinished:
		if msg.error != nil {
			log.Error("failed to notify that the work is finished: %v", msg.error)
		}

		if err := m.repository.CreateWorkTask(
			work.Task{
				Name:     m.task,
				Duration: m.choice,
			},
		); err != nil {
			log.Error("failed to save work task in the database", err)
		}

		return m, tea.Quit
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.isQuitting = true

			return m, tea.Quit
		case "enter":
			switch m.state {
			case taskSelection:
				return m, selectTask(m.taskInput.Value())
			case timeSelection:
				selectedItem := m.timeList.SelectedItem().(listItem)
				if selectedItem.label == "Custom Value" {
					return m, selectTimeFromInput(selectedItem.value)
				}

				return m, selectTimeFromList(selectedItem.value)
			case customTimeSelection:
				return m, selectCustomTime(m.timeInput.Value())
			case progressView:
				return m, tick()
			}
		}
	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, finishWork(m.choice)
		}

		m.timeRemaining--

		increment := 1.0 / float64(m.choice)
		cmd = m.progress.IncrPercent(increment)

		return m, tea.Batch(tick(), cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)

		return m, cmd
	}

	m.taskInput, cmd = m.taskInput.Update(msg)
	commands = append(commands, cmd)

	m.timeList, cmd = m.timeList.Update(msg)
	commands = append(commands, cmd)

	m.timeInput, cmd = m.timeInput.Update(msg)
	commands = append(commands, cmd)

	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	if m.isQuitting {
		return infoTextStyle.Render("Not working? That’s cool. Enjoy a break!")
	}

	switch m.state {
	case taskSelection:
		return fmt.Sprintf(
			"\n    What task do you want to work on?\n\n    %s\n\n    %s",
			m.taskInput.View(),
			"(q to quit)",
		)
	case timeSelection:
		return "\n" + m.timeList.View()

	case customTimeSelection:
		return fmt.Sprintf(
			"\n    How many minutes do you want to work for?\n\n    %s\n\n    %s",
			m.timeInput.View(),
			"(q to quit)",
		)
	case progressView:
		pad := strings.Repeat(" ", progressPadding)
		return infoTextStyle.Render(
			fmt.Sprintf(
				"Running timer for %d seconds. Have fun! %d seconds remaining.",
				m.choice,
				m.timeRemaining,
			),
		) + "\n" + pad + m.progress.View() + "\n\n" + pad + progressHelpStyle("Press q key to quit")
	default:
		return ""
	}
}
