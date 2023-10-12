package wshandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// CloseWS ... Closes Websocket connection

// NewConnection ... Connect to HS WSS Websocket

func (ctx *Ws) startWsConnection() {

	var cancel context.CancelFunc
	ctx.Con = context.Background()
	defer cancel()
	var err error
	ctx.C, _, err = websocket.Dial(ctx.Con, ctx.URI, nil)
	if err != nil {
		// ...
	}
	defer ctx.C.Close(websocket.StatusInternalError, "the sky is falling")

	go ctx.wsReader()

	ctx.wg.Done()

	<-ctx.CloseWS
	ctx.C.Close(websocket.StatusNormalClosure, "Client going offline.")
}

// NewConnection todo
func NewConnection(URI string) *Ws {

	var ctx Ws
	ctx.URI = URI
	ctx.wg.Add(1)
	go ctx.startWsConnection()
	ctx.wg.Wait()
	return &ctx

}

// QueryByID todo
func (ctx *Ws) QueryByID(Type string, ID string) {

	newLookup := LookupData{
		ID:   ID,
		Type: Type,
	}

	newMessage := WsMessage{
		Type: ctx.LookupType(),
		Data: newLookup,
		Pid:  uuid.NewString(),
	}

	wsjson.Write(ctx.Con, ctx.C, newMessage)

}

// Subscribe todo
func (ctx *Ws) Subscribe(Type []string) {

	newMessage := WsMessage{
		Type: ctx.subscribeType(),
		Data: ctx.Subscriptions.All(),
		Pid:  uuid.NewString(),
	}

	err := wsjson.Write(ctx.Con, ctx.C, newMessage)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// WriteMessage todo
func (ctx *Ws) WriteMessage(v any) {

}

func (ctx *Ws) wsReader() {

	for {
		_, wsMsgByte, err := ctx.C.Reader(ctx.Con)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		buffer := make([]byte, 1024)
		var message string
		for {
			count, err := wsMsgByte.Read(buffer)
			if count > 0 {
				message = string(buffer[:count])
				if ctx.last == message {
					break
				}
				ctx.last = message

				_, err = json.Marshal(message)
				if err == nil {
					var pong = PongMessage{
						Type: "pong",
						PID:  uuid.NewString(),
					}
					err = wsjson.Write(ctx.Con, ctx.C, pong)
				}
			}
			if err != nil {
				if err == io.EOF {
					fmt.Println(message)
					break
				} else {
					fmt.Println(err.Error())
				}
			}
		}

	}
}
