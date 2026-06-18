package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/recall/internal/repository"
)

func (app *App) SaveMemory(ctx context.Context, numSessions int, numCommands int) error {
	scanner := bufio.NewScanner(os.Stdin)

	// 1. List Sessions
	sessions, err := app.Store.Sessions.ListRecent(ctx, repository.Page{Limit: int64(numSessions)})
	if err != nil {
		return fmt.Errorf("failed to list sessions: %w", err)
	}

	if len(sessions) == 0 {
		fmt.Println("No recent sessions found. Try running some commands first.")
		return nil
	}

	fmt.Println("\nRecent Sessions:")
	for i, s := range sessions {
		startTime := time.UnixMilli(s.StartTs).Format("2006-01-02 15:04:05")
		fmt.Printf("[%d] %s (%s) - %d commands\n", i+1, s.Repo, startTime, s.CommandCount)
	}

	fmt.Print("\nSelect a session by number: ")
	if !scanner.Scan() {
		return nil
	}
	choice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || choice < 1 || choice > len(sessions) {
		return fmt.Errorf("invalid selection")
	}

	selectedSession := sessions[choice-1]

	// 2. Show Commands
	commands, err := app.Store.Commands.ListBySession(ctx, selectedSession.ID, repository.Page{Limit: int64(numCommands)})
	if err != nil {
		return fmt.Errorf("failed to list commands for session: %w", err)
	}

	fmt.Printf("\nCommands in session [%s]:\n", selectedSession.Repo)
	for _, c := range commands {
		fmt.Printf("  - %s\n", c.Command)
	}

	// 3. Gather Summary
	fmt.Print("\nWhat was the problem about?\n> ")
	if !scanner.Scan() {
		return nil
	}
	problem := strings.TrimSpace(scanner.Text())

	fmt.Print("\nWhat did you do to fix this?\n> ")
	if !scanner.Scan() {
		return nil
	}
	fix := strings.TrimSpace(scanner.Text())

	if problem == "" || fix == "" {
		fmt.Println("Problem and fix cannot be empty. Aborting.")
		return nil
	}

	// Format for FTS5 readiness
	summary := fmt.Sprintf("Problem:\n%s\n\nFix:\n%s", problem, fix)
	title := fmt.Sprintf("Memory for %s", selectedSession.Repo)

	// 4. Save Memory
	_, err = app.Store.Memories.Create(ctx, repository.CreateMemoryParams{
		ID:        uuid.NewString(),
		SessionID: selectedSession.ID,
		Title:     &title,
		Summary:   summary,
		CreatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		return fmt.Errorf("failed to save memory: %w", err)
	}

	fmt.Println("\nMemory saved successfully!")
	return nil
}
