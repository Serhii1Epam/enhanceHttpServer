package wsserver

import (
	"log"
	"net/http"

	"github.com/Serhii1Epam/enhanceHttpServer/pkg/hub"
	"github.com/gorilla/websocket"
)

type WsData struct {
	wsConn   *websocket.Conn
	upgrader *websocket.Upgrader
}

func NewWsData() *WsData {
	return &WsData{upgrader: &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	},
	}
}

func (wd *WsData) StartWs(h *hub.Hub, w http.ResponseWriter, r *http.Request) error {
	var err error
	wd.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	wd.wsConn, err = wd.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	client := &hub.Client{Hub: h, Conn: wd.wsConn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WriteMsg()
	go client.ReadMsg()

	return nil
}
