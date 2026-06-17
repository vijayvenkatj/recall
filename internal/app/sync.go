package app

import (
	"context"
	"io"
)


type Event struct {
	Timestamp int64
	ExitCode  int
	CWD       string
	Repo      string
	Command   string
}


func (app *App) Sync(ctx context.Context, ) {
}






func parseEvent(event io.Reader) {
	
}