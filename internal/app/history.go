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

// Reuses searchState / stateList / stateDetail and the shared header/section/
// footer styles defined in search.go (same package).

var (
	exitOKStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#3FB950")).Bold(true)
	exitFailStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F85149")).Bold(true)
)

type historyItem struct {
	session repository.Session
}

func (i historyItem) Title() string {
	return i.session.Repo
}

func (i historyItem) Description() string {
	date := time.UnixMilli(i.session.EndTs).Format("2006-01-02 15:04")
	return fmt.Sprintf("%s • %d commands", date, i.session.CommandCount)
}

func (i historyItem) FilterValue() string {
	date := time.UnixMilli(i.session.EndTs).Format("2006-01-02")
	return i.session.Repo + " " + date
}

type historyModel struct {
	app      *App
	ctx      context.Context
	list     list.Model
	viewport viewport.Model
	state    searchState
	width    int
	height   int
	err      error
}

func (m historyModel) Init() tea.Cmd {
	return nil
}

func (m historyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == stateList {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "q":
				if m.list.FilterState() != list.Filtering {
					return m, tea.Quit
				}
			case "enter":
				selected, ok := m.list.SelectedItem().(historyItem)
				if !ok {
					return m, nil
				}
				cmdList, err := m.app.Store.Commands.ListBySession(m.ctx, selected.session.ID, repository.Page{Limit: 500})
				if err != nil {
					m.err = err
					return m, nil
				}
				m.state = stateDetail
				m.viewport.SetContent(m.renderSession(selected.session, cmdList))
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
		m.viewport.Height = msg.Height - v - 6
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

func exitMark(c repository.Command) string {
	if !c.ExitCode.Valid {
		return " "
	}
	if c.ExitCode.Int64 == 0 {
		return exitOKStyle.Render("✓")
	}
	return exitFailStyle.Render("✗")
}

func (m historyModel) renderSession(sess repository.Session, cmds []repository.Command) string {
	var s strings.Builder

	title := fmt.Sprintf("%s  —  %s", sess.Repo, time.UnixMilli(sess.EndTs).Format("2006-01-02 15:04"))
	s.WriteString(lipgloss.NewStyle().
		Width(m.viewport.Width - 2).
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Render(title))
	s.WriteString("\n\n")

	s.WriteString(sectionHeaderStyle.Render(" COMMAND HISTORY "))
	s.WriteString("\n\n")

	// Commands arrive newest-first; show them chronologically like a transcript.
	for i := len(cmds) - 1; i >= 0; i-- {
		c := cmds[i]
		line := fmt.Sprintf("%s %s", exitMark(c), c.Command)
		s.WriteString(lipgloss.NewStyle().
			Width(m.viewport.Width - 4).
			Render(line))
		s.WriteString("\n")
	}

	return s.String()
}

func (m historyModel) View() string {
	if m.err != nil {
		return lipgloss.NewStyle().Margin(1, 2).Render(fmt.Sprintf("Error: %v\nPress q to quit", m.err))
	}

	if m.state == stateList {
		return lipgloss.NewStyle().Margin(1, 2).Render(m.list.View())
	}

	header := headerStyle.Render("SESSION COMMANDS")
	footer := footerStyle.Render(" ↑/↓: scroll • esc: back • q: quit")

	return lipgloss.NewStyle().Margin(1, 2).Render(
		fmt.Sprintf("%s\n\n%s\n%s", header, m.viewport.View(), footer),
	)
}

func (app *App) History(ctx context.Context, query string, limit int) error {
	cleanQuery := strings.TrimSpace(query)

	var sessions []repository.Session
	var err error

	if cleanQuery == "" {
		sessions, err = app.Store.Sessions.ListRecent(ctx, repository.Page{Limit: int64(limit)})
	} else {
		sessions, err = app.Store.Sessions.SearchByCommand(ctx, cleanQuery, int64(limit))
	}

	if err != nil {
		return fmt.Errorf("history failed: %w", err)
	}

	if len(sessions) == 0 {
		if cleanQuery == "" {
			fmt.Println("No command history yet. Your shell hook records commands as you work — run some commands, then try 'recall history' again.")
		} else {
			fmt.Printf("No sessions found containing: '%s'\n", cleanQuery)
		}
		return nil
	}

	items := make([]list.Item, len(sessions))
	for i, s := range sessions {
		items[i] = historyItem{session: s}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	if cleanQuery == "" {
		l.Title = "RECENT SESSIONS"
	} else {
		l.Title = fmt.Sprintf("SESSIONS MATCHING: %s", strings.ToUpper(cleanQuery))
	}
	l.Styles.Title = headerStyle

	m := historyModel{
		app:   app,
		ctx:   ctx,
		list:  l,
		state: stateList,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
