package app

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/recall/internal/repository"
)

type Event struct {
	Timestamp int64
	ExitCode  int
	CWD       string
	Repo      string
	Command   string
}

func (app *App) Sync(ctx context.Context) error {

	offset, err := app.getLastOffset(ctx)
	if err != nil {
		return err
	}

	newOffset, count, err := app.readLogs(ctx, offset)
	if err != nil {
		return err
	}

	_, err = app.Store.Metadata.Set(ctx, "last_offset", strconv.FormatInt(newOffset, 10))
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Printf("Synced %d new command(s) into database.\n", count)
	}

	return nil
}

func (app *App) readLogs(ctx context.Context, offset int64) (int64, int, error) {
	file, err := os.Open(app.Config.LogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	defer file.Close()

	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		return offset, 0, err
	}

	reader := bufio.NewReader(file)
	currentOffset := offset
	count := 0

	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			currentOffset += int64(len(line))
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return offset, 0, err
		}

		event, err := parseEvent(line)
		if err != nil {
			continue
		}

		if err := app.processEvent(ctx, event); err != nil {
			return offset, 0, err
		}
		count++
	}

	return currentOffset, count, nil
}

func (app *App) processEvent(ctx context.Context, event Event) error {
	session, err := sessionForEvent(ctx, event, app.Store.Sessions)
	if err != nil {
		return err
	}

	exitCode := int64(event.ExitCode)
	_, err = app.Store.Commands.Append(ctx, repository.CreateCommandParams{
		ID:        uuid.NewString(),
		SessionID: session.ID,

		Command: event.Command,
		CWD:     &event.CWD,
		Repo:    &event.Repo,

		ExitCode:  &exitCode,
		Timestamp: event.Timestamp,
		CreatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		return err
	}

	return nil
}

func parseEvent(logline string) (Event, error) {
	var event Event

	info := strings.SplitN(logline, "\t", 5)
	if len(info) != 5 {
		return event, fmt.Errorf("invalid event format")
	}

	timestamp, err := strconv.ParseInt(info[0], 10, 64)
	if err != nil {
		return event, err
	}

	exitCode, err := strconv.Atoi(info[1])
	if err != nil {
		return event, err
	}

	event.Timestamp = timestamp * 1000
	event.ExitCode = exitCode
	event.CWD = info[2]
	event.Repo = info[3]
	event.Command = strings.ReplaceAll(strings.TrimSpace(info[4]), "\\n", "\n")

	return event, nil
}

func sessionForEvent(ctx context.Context, event Event, sessionRepo *repository.SessionRepository) (repository.Session, error) {

	var session repository.Session

	sessions, err := sessionRepo.ListRecent(ctx, repository.Page{Limit: 1})
	if err != nil {
		return session, err
	}
	if len(sessions) == 0 {
		session, err = sessionRepo.Create(ctx, repository.CreateSessionParams{
			ID:           uuid.NewString(),
			Repo:         event.Repo,
			StartTs:      event.Timestamp,
			EndTs:        event.Timestamp,
			CommandCount: 0, // Starts at 0; touched in Append to become 1
			CreatedAt:    time.Now().UnixMilli(),
			UpdatedAt:    time.Now().UnixMilli(),
		})
		if err != nil {
			return session, err
		}
		return session, nil
	}

	latestSession := sessions[0]
	last := time.UnixMilli(latestSession.EndTs)
	current := time.UnixMilli(event.Timestamp)

	if current.Sub(last) > 30*time.Minute || latestSession.Repo != event.Repo {
		session, err = sessionRepo.Create(ctx, repository.CreateSessionParams{
			ID:           uuid.NewString(),
			Repo:         event.Repo,
			StartTs:      event.Timestamp,
			EndTs:        event.Timestamp,
			CommandCount: 0, // Starts at 0; touched in Append to become 1
			CreatedAt:    time.Now().UnixMilli(),
			UpdatedAt:    time.Now().UnixMilli(),
		})
		if err != nil {
			return session, err
		}
		return session, nil
	}

	// Just return latestSession; it will be touched by Append inside the transaction
	return latestSession, nil
}

func (app *App) getLastOffset(ctx context.Context) (int64, error) {
	meta, err := app.Store.Metadata.Get(ctx, "last_offset")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return strconv.ParseInt(meta.Value, 10, 64)
}
