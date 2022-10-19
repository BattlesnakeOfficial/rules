package commands

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/BattlesnakeOfficial/rules/client"
)

type GameExporter struct {
	game          client.Game
	snakeRequests []client.SnakeRequest
	winner        SnakeState
	isDraw        bool
}

type result struct {
	WinnerID   string `json:"winnerId"`
	WinnerName string `json:"winnerName"`
	IsDraw     bool   `json:"isDraw"`
}

func (ge *GameExporter) FlushToFile(outputFile io.Writer) (int, error) {
	formattedOutput, err := ge.ConvertToJSON()
	if err != nil {
		return 0, err
	}

	for _, line := range formattedOutput {
		_, err := io.WriteString(outputFile, fmt.Sprintf("%s\n", line))
		if err != nil {
			return 0, err
		}
	}

	return len(formattedOutput), nil
}

func (ge *GameExporter) ConvertToJSON() ([]string, error) {
	output := make([]string, 0)
	serialisedGame, err := json.Marshal(ge.game)
	if err != nil {
		return output, err
	}
	output = append(output, string(serialisedGame))
	for _, board := range ge.snakeRequests {
		serialisedBoard, err := json.Marshal(board)
		if err != nil {
			return output, err
		}
		output = append(output, string(serialisedBoard))
	}
	serialisedResult, err := json.Marshal(result{
		WinnerID:   ge.winner.ID,
		WinnerName: ge.winner.Name,
		IsDraw:     ge.isDraw,
	})
	if err != nil {
		return output, err
	}
	output = append(output, string(serialisedResult))
	return output, nil
}

func (ge *GameExporter) AddSnakeRequest(snakeRequest client.SnakeRequest) {
	ge.snakeRequests = append(ge.snakeRequests, snakeRequest)
}
