package commands

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type CLIWebsocketServer struct {
	game   Game
	events chan GameEvent
	done   chan bool
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewCLIWebsocketServer(game Game) *CLIWebsocketServer {
	return &CLIWebsocketServer{
		game:   game,
		events: make(chan GameEvent, 1000),
		done:   make(chan bool),
	}
}

func (server *CLIWebsocketServer) handleGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(struct {
		Game Game
	}{server.game})
	if err != nil {
		log.Printf("Unable to serialize game for /games/: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (server *CLIWebsocketServer) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Unable to upgrade connection: %v", err)
		return
	}

	defer func() {
		err = ws.Close()
		if err != nil {
			log.Printf("Unable to close websocket stream")
		}
	}()

	for event := range server.events {
		jsonStr, err := json.Marshal(event)
		if err != nil {
			log.Printf("Unable to serialize event for websocket: %v", err)
		}

		err = ws.WriteMessage(websocket.TextMessage, jsonStr)
		if err != nil {
			log.Printf("Unable to write to websocket: %v", err)
			break
		}
	}

	log.Printf("Finished writing all game events, signalling game server to stop")
	close(server.done)

	log.Printf("Sending websocket close message")
	err = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Printf("Problem closing websocket: %v", err)
	}
}

func (server *CLIWebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/games/") {
		if strings.HasSuffix(r.URL.Path, "/events") {
			server.handleWebsocket(w, r)
			return
		}
		server.handleGame(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (server *CLIWebsocketServer) listen() *httptest.Server {
	return httptest.NewServer(
		cors.Default().Handler(server),
	)
}
