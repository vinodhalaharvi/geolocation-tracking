package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

type WsHandler struct {
}

//func main() {
//	r := gin.Default()
//	r.LoadHTMLFiles("templates/map.html")
//
//	r.GET("/", func(c *gin.Context) {
//		c.HTML(http.StatusOK, "map.html", nil)
//	})
//
//	r.GET("/ws", func(c *gin.Context) {
//		Handle(c.Writer, c.Request)
//	})
//
//	r.Run(":8080") // Listen and serve on 0.0.0.0:8080
//}

func (ws *WsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("error: %v", err)
		}
	}(conn)

	for {
		// Simulate marker updates
		markerUpdate := map[string]float64{
			"latitude":  37.4220656,
			"longitude": -122.0840897,
		}

		if err := conn.WriteJSON(markerUpdate); err != nil {
			log.Println(err)
			break
		}

		time.Sleep(2 * time.Second) // Send updates every 2 seconds
	}
}
