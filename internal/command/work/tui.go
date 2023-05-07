package work

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/sergiught/work-pilot-cli/internal/work"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/gen2brain/beeep"
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

type Model struct {
	repository *work.Repository

	taskInput textinput.Model
	timeList  list.Model
	timeInput textinput.Model
	progress  progress.Model

	isQuitting         bool
	taskSelected       bool
	listTimeSelected   bool
	customTimeSelected bool

	choice        int
	timeRemaining int
	task          string
}

func NewWorkModel(repository *work.Repository) *Model {
	customTaskInput := textinput.New()
	customTaskInput.SetValue("Work")
	customTaskInput.Placeholder = "Work"
	customTaskInput.Focus()
	customTaskInput.CharLimit = 256
	customTaskInput.Width = 256

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

	customTimeInput := textinput.New()
	customTimeInput.Placeholder = "0"
	customTimeInput.Focus()
	customTimeInput.CharLimit = 32
	customTimeInput.Width = 20

	progressIndicator := progress.New(progress.WithDefaultGradient())

	return &Model{
		repository: repository,
		taskInput:  customTaskInput,
		timeList:   timeList,
		timeInput:  customTimeInput,
		progress:   progressIndicator,
	}
}

func (m Model) Init() tea.Cmd {
	if m.choice != 0 {
		return tickCmd()
	}

	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.timeList.SetWidth(msg.Width)

		m.progress.Width = msg.Width - progressPadding*2 - 4
		if m.progress.Width > progressMaxWidth {
			m.progress.Width = progressMaxWidth
		}

		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.isQuitting = true

			return m, tea.Quit
		case "enter":
			if m.taskSelected && (m.listTimeSelected || m.customTimeSelected) && m.choice != 0 {
				return m, tickCmd()
			}

			if m.taskSelected && m.listTimeSelected && m.customTimeSelected && m.choice == 0 {
				value, err := strconv.Atoi(m.timeInput.Value())
				if err != nil {
					log.Error("failed to convert value to int", "value", value, "err", err)
					return m, tea.Quit
				}

				m.choice = value
				m.timeRemaining = value

				return m, tickCmd()
			}

			if m.taskSelected && !m.listTimeSelected {
				if selectedItem, ok := m.timeList.SelectedItem().(listItem); ok {
					m.listTimeSelected = true

					if selectedItem.label == "Custom Value" {
						m.customTimeSelected = true

						m.timeInput.Reset()

						var cmd tea.Cmd
						m.timeInput, cmd = m.timeInput.Update(msg)
						return m, cmd
					}

					m.choice = selectedItem.value
					m.timeRemaining = selectedItem.value

					return m, tickCmd()
				}
			}

			if m.taskInput.Value() != "" && (!m.listTimeSelected || !m.customTimeSelected) {
				m.taskSelected = true
				m.task = m.taskInput.Value()

				var cmd tea.Cmd
				m.timeList, cmd = m.timeList.Update(msg)
				return m, cmd
			}

			var cmd tea.Cmd
			m.timeList, cmd = m.timeList.Update(msg)
			return m, cmd
		}
	case tickMsg:
		if m.progress.Percent() == 1.0 {
			if err := beeep.Beep(44000, 10000); err != nil {
				log.Error("failed to notify with a beep that work finished", err)
			}

			if err := beeep.Notify(
				"Work Pilot: Work Finished!",
				fmt.Sprintf("Congratulations! You've worked for %d second(s).", m.choice),
				"",
			); err != nil {
				log.Error("failed to notify with a notification that work finished", err)
			}

			task := work.Task{
				Name:     m.task,
				Duration: m.choice,
			}

			if err := m.repository.CreateWorkTask(task); err != nil {
				log.Error("failed to save work task in the database", err)
			}

			return m, tea.Quit
		}

		m.timeRemaining--

		increment := 1.0 / float64(m.choice)
		cmd := m.progress.IncrPercent(increment)

		return m, tea.Batch(tickCmd(), cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)

		return m, cmd
	}

	var commands []tea.Cmd

	var cmd tea.Cmd
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

	if m.choice != 0 {
		pad := strings.Repeat(" ", progressPadding)
		return infoTextStyle.Render(
			fmt.Sprintf(
				"Running timer for %d seconds. Have fun! %d seconds remaining.",
				m.choice,
				m.timeRemaining,
			),
		) +
			"\n" +
			pad + m.progress.View() +
			"\n\n" +
			pad + progressHelpStyle("Press q key to quit")
	}

	if !m.taskSelected {
		return "\n    " +
			fmt.Sprintf(
				"What task do you want to work on?\n\n    %s\n\n    %s",
				m.taskInput.View(),
				"(q to quit)",
			) +
			"\n"
	}

	if m.taskSelected && m.listTimeSelected && m.customTimeSelected {
		return "\n    " +
			fmt.Sprintf(
				"How many minutes do you want to work for?\n\n    %s\n\n    %s",
				m.timeInput.View(),
				"(q to quit)",
			) +
			"\n"
	}

	if m.taskSelected && (!m.listTimeSelected || !m.customTimeSelected) {
		return "\n" + m.timeList.View()
	}

	return ""
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type listItem struct {
	label string
	value int
}

func (i listItem) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
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
