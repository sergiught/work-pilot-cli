package work

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	progressHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
)

type Model struct {
	list       list.Model
	progress   progress.Model
	isQuitting bool

	Choice int
}

func NewWorkModel() Model {
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
	}

	l := list.New(items, itemDelegate{}, listWidth, listHeight)
	l.Title = "How many minutes do you want to work for?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	p := progress.New(progress.WithDefaultGradient())

	return Model{
		list:     l,
		progress: p,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.Choice != 0 {
		pad := strings.Repeat(" ", progressPadding)
		return quitTextStyle.Render(fmt.Sprintf("Running timer for %d seconds. Have fun!", m.Choice)) + "\n" + pad + m.progress.View() + "\n\n" + pad + progressHelpStyle("Press any key to quit")
	}
	if m.isQuitting {
		return quitTextStyle.Render("Not working? That’s cool. Enjoy a break!")
	}
	return "\n" + m.list.View()
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
			i, ok := m.list.SelectedItem().(listItem)
			if ok {
				m.Choice = i.value
				return m, tickCmd()
			}
			return m, tea.Quit
		}
	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		increment := 1.0 / float64(m.Choice)
		cmd := m.progress.IncrPercent(increment)
		return m, tea.Batch(tickCmd(), cmd)
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	default:
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
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

	str := fmt.Sprintf("%d. %s", index+1, i.label)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
