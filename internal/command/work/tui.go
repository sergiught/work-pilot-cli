package work

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

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
	list     list.Model
	input    textinput.Model
	progress progress.Model

	isQuitting      bool
	inputIsSelected bool
	choice          int
	timeRemaining   int
}

func NewWorkModel() *Model {
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

	l := list.New(items, itemDelegate{}, listWidth, listHeight)
	l.Title = "How many minutes do you want to work for?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	ti := textinput.New()
	ti.Placeholder = "0"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 20

	p := progress.New(progress.WithDefaultGradient())

	return &Model{
		list:     l,
		input:    ti,
		progress: p,
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
		m.list.SetWidth(msg.Width)

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
			if m.choice != 0 {
				return m, tickCmd()
			}

			if m.inputIsSelected {
				value, err := strconv.Atoi(m.input.Value())
				if err != nil {
					return m, tea.Quit
				}
				m.choice = value
				m.timeRemaining = value
				return m, tickCmd()
			}

			i, ok := m.list.SelectedItem().(listItem)
			if ok {
				if i.label == "Custom Value" {
					m.inputIsSelected = true

					var commands []tea.Cmd

					var cmd tea.Cmd
					m.list, cmd = m.list.Update(msg)
					commands = append(commands, cmd)

					m.input, cmd = m.input.Update(msg)
					commands = append(commands, cmd)

					return m, tea.Batch(commands...)
				}

				m.choice = i.value
				m.timeRemaining = i.value

				return m, tickCmd()
			}

			return m, tea.Quit
		}
	case tickMsg:
		if m.progress.Percent() == 1.0 {
			err := beeep.Beep(44000, 10000)
			if err != nil {
				log.Error("failed to notify with a beep that work finished", err)
			}

			err = beeep.Notify(
				"Work Pilot: Work Finished!",
				fmt.Sprintf("Congratulations! You've worked for %d second(s).", m.choice),
				"",
			)
			if err != nil {
				log.Error("failed to notify with a notification that work finished", err)
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
	m.list, cmd = m.list.Update(msg)
	commands = append(commands, cmd)

	m.input, cmd = m.input.Update(msg)
	commands = append(commands, cmd)

	return m, tea.Batch(commands...)
}

func (m Model) View() string {
	if m.isQuitting {
		return infoTextStyle.Render("Not working? That’s cool. Enjoy a break!")
	}

	if m.inputIsSelected && m.choice == 0 {
		return "\n    " +
			fmt.Sprintf(
				"How many minutes do you want to work for?\n\n    %s\n\n    %s",
				m.input.View(),
				"(q to quit)",
			) +
			"\n"
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

	return "\n" + m.list.View()
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
