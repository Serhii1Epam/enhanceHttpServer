package wsserver

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 1024
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

func (wd *WsData) StartWs(w http.ResponseWriter, r *http.Request) error {
	var err error
	wd.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	wd.wsConn, err = wd.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("1:FAILED.\n")
		fmt.Printf("%v\n", wd.wsConn)
		log.Println(err)
		return err
	}

	wd.OpenWs()

	/*if err = wd.openWs(); err != nil {
		log.Printf("2:FAILED.\n")
		log.Println(err)
		return err
	}*/

	return nil
}

func (wd *WsData) OpenWs() error {
	if wd.wsConn == nil {
		return errors.New("Connection not created.")
	}

	defer wd.wsConn.Close()

	wd.wsConn.SetReadLimit(maxMessageSize)
	wd.wsConn.SetReadDeadline(time.Now().Add(pongWait))
	wd.wsConn.SetPongHandler(func(string) error { wd.wsConn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		msgType, msgData, err := wd.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msgData = bytes.TrimSpace(bytes.Replace(msgData, []byte{'\n'}, []byte{' '}, -1))
		fmt.Printf("msgData:[%v]", msgData)

		if err := wd.wsConn.WriteMessage(msgType, msgData); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
