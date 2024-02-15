package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vinodhalaharvi/geolocation-tracking/service"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Websocket struct {
	Clients map[*websocket.Conn]bool
	Lock    sync.Mutex
}

func NewWebsocket() *Websocket {
	return &Websocket{
		Clients: make(map[*websocket.Conn]bool),
		Lock:    sync.Mutex{},
	}
}

func (ws *Websocket) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Add the new connection to the Clients map
	ws.Lock.Lock()
	ws.Clients[conn] = true
	ws.Lock.Unlock()

	// Goroutine to manage the connection's lifecycle
	go func() {
		defer func() {
			ws.Lock.Lock()
			delete(ws.Clients, conn) // Clean up after disconnection
			ws.Lock.Unlock()

			err := conn.Close() // Ensure the connection is closed
			if err != nil {
				log.Printf("Failed to close connection: %v", err)
			}
		}()

		for {
			_, _, err := conn.NextReader()
			if err != nil {
				log.Println("Read error:", err)
				break // Exit the loop on error, triggering cleanup
			}
		}
	}()
}

func (ws *Websocket) BroadcastMessage(message string) {
	ws.Lock.Lock()
	defer ws.Lock.Unlock()

	for client := range ws.Clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			// Remove the client from the map if sending the message fails
			delete(ws.Clients, client)
			err := client.Close()
			if err != nil {
				return
			}
		}
	}
}

func (ws *Websocket) RunSimulation(c *gin.Context, gfs *service.GeoFenceStateService) {
	ticker := time.NewTicker(1 * time.Second) // Create a ticker for 1-second intervals
	defer ticker.Stop()                       // Ensure the ticker is stopped to free resources

	for {
		select {
		case <-ticker.C:
			// It's time to send another update
			message, err := ws.CreateUpdateMessage(gfs)
			if err != nil {
				log.Printf("Error creating update message: %v", err)
				continue
			}
			ws.BroadcastMessage(message)
		}
	}
}
func (ws *Websocket) CreateUpdateMessage(gfs *service.GeoFenceStateService) (string, error) {
	// Call the method to simulate car movements before creating the message
	gfs.SimulateAssetMovements()

	jsonData, err := json.MarshalIndent(gfs.State.GetAssets(), "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
