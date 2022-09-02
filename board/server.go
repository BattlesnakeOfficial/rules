package board

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	log "github.com/spf13/jwalterweatherman"
)

// A minimal server capable of handling the requests from a single browser client running the board viewer.
type BoardServer struct {
	game   Game
	events chan GameEvent // channel for sending events from the game runner to the browser client
	done   chan bool      // channel for signalling (via closing) that all events have been sent to the browser client

	httpServer *http.Server
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewBoardServer(game Game) *BoardServer {
	mux := http.NewServeMux()

	server := &BoardServer{
		game:   game,
		events: make(chan GameEvent, 1000), // buffered channel to allow game to run ahead of browser client
		done:   make(chan bool),
		httpServer: &http.Server{
			Handler: cors.Default().Handler(mux),
		},
	}

	mux.HandleFunc("/games/"+game.ID, server.handleGame)
	mux.HandleFunc("/games/"+game.ID+"/events", server.handleWebsocket)

	return server
}

// Handle the /games/:id request made by the board to fetch the game metadata.
func (server *BoardServer) handleGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(struct {
		Game Game
	}{server.game})
	if err != nil {
		log.ERROR.Printf("Unable to serialize game: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Handle the /games/:id/events websocket request made by the board to receive game events.
func (server *BoardServer) handleWebsocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.ERROR.Printf("Unable to upgrade connection: %v", err)
		return
	}

	defer func() {
		err = ws.Close()
		if err != nil {
			log.ERROR.Printf("Unable to close websocket stream")
		}
	}()

	for event := range server.events {
		jsonStr, err := json.Marshal(event)
		if err != nil {
			log.ERROR.Printf("Unable to serialize event for websocket: %v", err)
		}

		err = ws.WriteMessage(websocket.TextMessage, jsonStr)
		if err != nil {
			log.ERROR.Printf("Unable to write to websocket: %v", err)
			break
		}
	}

	log.DEBUG.Printf("Finished writing all game events, signalling game server to stop")
	close(server.done)

	log.DEBUG.Printf("Sending websocket close message")
	err = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.ERROR.Printf("Problem closing websocket: %v", err)
	}
}

func (server *BoardServer) Listen() (string, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	go func() {
		err = server.httpServer.Serve(listener)
		if err != http.ErrServerClosed {
			log.ERROR.Printf("Error in board HTTP server: %v", err)
		}
	}()

	url := "http://" + listener.Addr().String()

	return url, nil
}

func (server *BoardServer) Shutdown() {
	close(server.events)

	log.DEBUG.Printf("Waiting for websocket clients to finish")
	<-server.done
	log.DEBUG.Printf("Server is done, exiting")

	err := server.httpServer.Shutdown(context.Background())
	if err != nil {
		log.ERROR.Printf("Error shutting down HTTP server: %v", err)
	}
}

func (server *BoardServer) SendEvent(event GameEvent) {
	server.events <- event
}
