package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/vijayvenkatj/recall/internal/llm"
	"github.com/vijayvenkatj/recall/internal/repository"
)

type saveStep int

const (
	stepSelectSession saveStep = iota
	stepReviewCommands
	stepInputProblem
	stepInputFix
	stepSaving
	stepDone
)

type saveModel struct {
	app            *App
	ctx            context.Context
	sessions       []repository.Session
	commands       []repository.Command
	viewport       viewport.Model
	problemInput   textinput.Model
	fixInput       textinput.Model
	selectedIdx    int
	step           saveStep
	err            error
	numCmds        int
	width          int
	height         int
	llmSuggestions llm.Suggestions
	llmLoading     bool
	llmError       error
	customTitle    string
}

func (m saveModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m saveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.step == stepSelectSession || m.step == stepReviewCommands || m.step == stepDone {
				return m, tea.Quit
			}
		case "esc":
			if m.step == stepSelectSession || m.step == stepDone {
				return m, tea.Quit
			}
			if m.step == stepReviewCommands {
				m.step = stepSelectSession
				return m, nil
			}
			if m.step == stepInputProblem {
				m.step = stepReviewCommands
				return m, nil
			}
			if m.step == stepInputFix {
				m.step = stepInputProblem
				m.problemInput.Focus()
				return m, nil
			}
		}

		switch m.step {
		case stepSelectSession:
			switch msg.String() {
			case "up", "k":
				if m.selectedIdx > 0 {
					m.selectedIdx--
				}
			case "down", "j":
				if m.selectedIdx < len(m.sessions)-1 {
					m.selectedIdx++
				}
			case "enter":
				m.step = stepReviewCommands
				// Fetch commands for review - fetch more to ensure "all" are seen
				limit := int64(m.numCmds)
				if limit < 100 {
					limit = 500
				}
				cmds, err := m.app.Store.Commands.ListBySession(m.ctx, m.sessions[m.selectedIdx].ID, repository.Page{Limit: limit})
				if err != nil {
					m.err = err
					return m, nil
				}
				m.commands = cmds
				
				// Initialize viewport
				var content strings.Builder
				for _, c := range cmds {
					content.WriteString(fmt.Sprintf("• %s\n", c.Command))
				}
				m.viewport.SetContent(CommandListStyle.Render(content.String()))
				m.viewport.YOffset = 0

				m.problemInput.SetValue("")
				m.fixInput.SetValue("")
			}
		case stepReviewCommands:
			if msg.String() == "enter" {
				m.step = stepInputProblem
				m.problemInput.Focus()
			} else {
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		case stepInputProblem:
			if msg.String() == "enter" && m.problemInput.Value() != "" {
				m.step = stepInputFix
				m.fixInput.Focus()
			} else {
				m.problemInput, cmd = m.problemInput.Update(msg)
				return m, cmd
			}
		case stepInputFix:
			if msg.String() == "enter" && m.fixInput.Value() != "" {
				m.step = stepSaving
				if m.app.LLMProvider != nil {
					m.llmLoading = true
					m.llmError = nil
					return m, m.generateSummaryCmd()
				} else {
					return m, m.saveManualMemory
				}
			} else {
				m.fixInput, cmd = m.fixInput.Update(msg)
				return m, cmd
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.viewport.Width = msg.Width - h
		m.viewport.Height = msg.Height - v - 6
		
		m.problemInput.Width = msg.Width - h - 4
		m.fixInput.Width = msg.Width - h - 4

	case llmSummaryMsg:
		m.llmLoading = false
		if msg.err != nil {
			m.llmError = msg.err
			// Gracefully fall back to manual values if LLM generation fails
			return m, m.saveManualMemory
		} else {
			return m, func() tea.Msg {
				err := m.saveMemoryWithValues(msg.result.Title, msg.result.Summary)
				return saveResultMsg{err: err}
			}
		}

	case saveResultMsg:
		if msg.err != nil {
			m.err = msg.err
			m.step = stepSelectSession
		} else {
			m.step = stepDone
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

type saveResultMsg struct{ err error }

type llmSummaryMsg struct {
	result llm.SummaryResult
	err    error
}

func (m saveModel) generateSummaryCmd() tea.Cmd {
	return func() tea.Msg {
		if m.app.LLMProvider == nil {
			return nil
		}

		var cmdList []string
		for _, c := range m.commands {
			cmdList = append(cmdList, c.Command)
		}

		result, err := m.app.LLMProvider.GenerateSummary(m.ctx, m.sessions[m.selectedIdx].Repo, cmdList, m.problemInput.Value(), m.fixInput.Value())
		return llmSummaryMsg{
			result: result,
			err:    err,
		}
	}
}

func (m saveModel) saveManualMemory() tea.Msg {
	summary := fmt.Sprintf("Problem:\n%s\n\nFix:\n%s", m.problemInput.Value(), m.fixInput.Value())
	title := fmt.Sprintf("Memory for %s", m.sessions[m.selectedIdx].Repo)
	err := m.saveMemoryWithValues(title, summary)
	return saveResultMsg{err: err}
}

func (m saveModel) saveMemoryWithValues(title string, summary string) error {
	_, err := m.app.Store.Memories.Create(m.ctx, repository.CreateMemoryParams{
		ID:        uuid.NewString(),
		SessionID: m.sessions[m.selectedIdx].ID,
		Title:     &title,
		Summary:   summary,
		CreatedAt: time.Now().UnixMilli(),
	})
	return err
}

func (m saveModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress q to quit", m.err)
	}

	var s strings.Builder
	s.WriteString(TitleStyle.Render(" SAVE MEMORY "))
	s.WriteString("\n\n")

	switch m.step {
	case stepSelectSession:
		s.WriteString("Select a session:\n\n")
		for i, sess := range m.sessions {
			cursor := " "
			style := lipgloss.NewStyle()
			if i == m.selectedIdx {
				cursor = ">"
				style = SelectedStyle
			}
			date := time.UnixMilli(sess.StartTs).Format("2006-01-02 15:04:05")
			s.WriteString(fmt.Sprintf("%s %s %s (%d cmds)\n", style.Render(cursor), style.Render(sess.Repo), SubtleStyle.Render(date), sess.CommandCount))
		}
		s.WriteString(SubtleStyle.Render("\n ↑/↓: navigate • enter: select • esc/q: quit"))
	case stepReviewCommands:
		s.WriteString(fmt.Sprintf("Reviewing commands for %s (↑/↓ to scroll):\n\n", m.sessions[m.selectedIdx].Repo))
		s.WriteString(m.viewport.View())
		s.WriteString(SubtleStyle.Render("\n\n enter: continue • esc: back • q: quit"))
	case stepInputProblem:
		s.WriteString("What was the problem about?\n\n")
		s.WriteString(m.problemInput.View())
		s.WriteString(SubtleStyle.Render("\n\n enter: next • esc: back"))
	case stepInputFix:
		s.WriteString("What did you do to fix this?\n\n")
		s.WriteString(m.fixInput.View())
		s.WriteString(SubtleStyle.Render("\n\n enter: save memory • esc: back"))
	case stepSaving:
		if m.llmLoading {
			s.WriteString(SubtleStyle.Render("🤖 [LLM: Summarizing session commands, problem statement, and fix...]") + "\n\n")
		} else {
			s.WriteString("Saving memory...")
		}
	case stepDone:
		s.WriteString("Done! Memory saved successfully.")
	}

	return lipgloss.NewStyle().Margin(1, 2).Render(s.String())
}

func (app *App) SaveMemory(ctx context.Context, numSessions int, numCommands int) error {
	sessions, err := app.Store.Sessions.ListRecent(ctx, repository.Page{Limit: int64(numSessions)})
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		fmt.Println("No recent sessions found.")
		return nil
	}

	pi := textinput.New()
	pi.Placeholder = "Describe the problem..."
	pi.CharLimit = 1000

	fi := textinput.New()
	fi.Placeholder = "Describe the fix..."
	fi.CharLimit = 1000

	m := saveModel{
		app:          app,
		ctx:          ctx,
		sessions:     sessions,
		numCmds:      numCommands,
		problemInput: pi,
		fixInput:     fi,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
