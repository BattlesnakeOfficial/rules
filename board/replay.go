package board

import (
	"os"
	"encoding/json"
	"fmt"
	"io"

	log "github.com/spf13/jwalterweatherman"
)

type ReplayFile struct {
	handle io.WriteCloser
}

func NewReplayFile(path string) *ReplayFile {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.ERROR.Fatalf("Failed to open replay file: %w", err)
	}

	return &ReplayFile{
		handle: fd,
	}
}

func (replay *ReplayFile) WriteGameInfo(game Game) {
	// TODO(schoon): Provide a clear delimiter between game info and frames.
	// Additionally, we probably want to ensure they're ordered. Take in the
	// `Game` in `NewReplayFile` (much like `NewBoardServer` and write
	// game info as front matter?
	jsonStr, err := json.Marshal(struct {
		Game Game
	}{game})
	if err != nil {
		log.ERROR.Printf("Unable to serialize event for replay file: %v", err)
	}

	_, err = io.WriteString(replay.handle, fmt.Sprintf("%s\n", jsonStr))
	if err != nil {
		log.WARN.Printf("Unable to write to replay file: %v", err)
	}
}

func (replay *ReplayFile) WriteEvent(event GameEvent) {
	jsonStr, err := json.Marshal(event)
	if err != nil {
		log.ERROR.Printf("Unable to serialize event for replay file: %v", err)
	}

	_, err = io.WriteString(replay.handle, fmt.Sprintf("%s\n", jsonStr))
	if err != nil {
		log.WARN.Printf("Unable to write to replay file: %v", err)
	}
}
