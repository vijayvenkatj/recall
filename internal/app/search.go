package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vijayvenkatj/recall/internal/repository"
)

type searchState int

const (
	stateList searchState = iota
	stateDetail
)

type item struct {
	memory   repository.Memory
	commands []repository.Command
}

func (i item) Title() string {
	if i.memory.Title.Valid {
		return i.memory.Title.String
	}
	return "Untitled Memory"
}

func (i item) Description() string {
	date := time.UnixMilli(i.memory.CreatedAt).Format("2006-01-02 15:04")
	return fmt.Sprintf("%s • %d commands", date, len(i.commands))
}

func (i item) FilterValue() string {
	val := i.Title() + " " + i.memory.Summary
	return val
}

type searchModel struct {
	list     list.Model
	viewport viewport.Model
	state    searchState
	width    int
	height   int
	query    string
}

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	detailTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true).
				MarginBottom(1)

	sectionHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#626262")).
				Padding(0, 1).
				MarginTop(1).
				MarginBottom(1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)
)

func (m searchModel) Init() tea.Cmd {
	return nil
}

func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == stateList {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				selected := m.list.SelectedItem().(item)
				m.state = stateDetail
				m.viewport.SetContent(m.renderDetail(selected))
				m.viewport.YOffset = 0
			}
		} else {
			switch msg.String() {
			case "esc", "backspace":
				m.state = stateList
			case "ctrl+c", "q":
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		m.viewport.Width = msg.Width - h
		m.viewport.Height = msg.Height - v - 6 // leave space for header/footer
	}

	if m.state == stateList {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m searchModel) renderDetail(i item) string {
	var s strings.Builder

	s.WriteString(detailTitleStyle.Render(i.Title()))
	s.WriteString("\n\n")
	s.WriteString(i.memory.Summary)
	s.WriteString("\n")

	if len(i.commands) > 0 {
		s.WriteString(sectionHeaderStyle.Render(" COMMAND HISTORY "))
		s.WriteString("\n")
		
		var cmds strings.Builder
		for _, c := range i.commands {
			cmds.WriteString(fmt.Sprintf("• %s\n", c.Command))
		}
		s.WriteString(CommandListStyle.Render(cmds.String()))
	}

	return s.String()
}

func (m searchModel) View() string {
	if m.state == stateList {
		return lipgloss.NewStyle().Margin(1, 2).Render(m.list.View())
	}

	header := headerStyle.Render("DETAIL VIEW")
	footer := footerStyle.Render(" ↑/↓: scroll • esc: back • q: quit")

	return lipgloss.NewStyle().Margin(1, 2).Render(
		fmt.Sprintf("%s\n\n%s\n%s", header, m.viewport.View(), footer),
	)
}

func (app *App) Search(ctx context.Context, query string) error {
	cleanQuery := strings.TrimSpace(query)
	if cleanQuery == "" {
		return nil
	}

	memories, err := app.Store.Memories.Search(ctx, cleanQuery, 20)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	items := make([]list.Item, len(memories))
	for i, m := range memories {
		cmds, _ := app.Store.Commands.ListBySession(ctx, m.SessionID, repository.Page{Limit: 500})
		items[i] = item{
			memory:   m,
			commands: cmds,
		}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = fmt.Sprintf("RESULTS FOR: %s", strings.ToUpper(cleanQuery))
	l.Styles.Title = headerStyle

	m := searchModel{
		list:  l,
		state: stateList,
		query: cleanQuery,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
