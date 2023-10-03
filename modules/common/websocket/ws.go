package wshandler

import (
	"encoding/json"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var WsWriteChan = make(chan WsWriteStruct)

var WsWrite = WsWriteMuWrapper{
	WsWrite: make(chan WsWriteStruct),
}

type WsWriteMuWrapper struct {
	WsWrite chan WsWriteStruct
	Mu      sync.Mutex
}

type WsWriteStruct struct {
	WsMsg     []byte
	WsMsgType int
}

var WsRead = make(chan []byte)

var ping struct {
	WsType string `json:"type"`
	Id     string `json:"pid"`
}
var pong struct {
	WsType string `json:"type"`
	Id     string `json:"pid"`
}

func WsConn() {

	u := url.URL{Scheme: "wss", Host: "hs.vtolvr.live", Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				panic(err)
			}

			Ping := ping

			err = json.Unmarshal(message, &Ping)
			if err == nil {
				Pong := pong
				Pong.Id = Ping.Id
				Pong.WsType = Ping.WsType
				PongByte, err := json.Marshal(Pong)
				if err != nil {
					panic(err)
				}

				WsWrite.Mu.Lock()
				WsWrite.WsWrite <- WsWriteStruct{WsMsg: PongByte, WsMsgType: websocket.PongMessage}

			}

			WsRead <- message
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	go func() {
		defer close(done)
		for {
			WsMsg := <-WsWrite.WsWrite
			err = c.WriteMessage(WsMsg.WsMsgType, WsMsg.WsMsg)
			if err != nil {
				panic(err)
			}
			WsWrite.Mu.Unlock()
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}

		}
	}
}
