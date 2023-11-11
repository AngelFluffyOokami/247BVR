package wshandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/angelfluffyookami/247BVR/modules/common/global"
	"github.com/angelfluffyookami/247BVR/modules/handlers"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// startWsConnection initiates a WebSocket connection
func (ctx *Ws) startWsConnection() {
	var cancel context.CancelFunc
	for {
		ctx.Con, cancel = context.WithCancel(context.Background())

		var err error
		ctx.C, _, err = websocket.Dial(ctx.Con, ctx.URI, nil)
		if err != nil {
			log.Printf("Failed to dial WebSocket: %v. Retrying...", err)
			cancel()
			time.Sleep(time.Duration(ctx.SecondPerAttempt) * time.Second) // Wait for a few seconds before retrying
			continue                                                      // Retry the connection
		}

		// Start the WebSocket reader Goroutine
		go ctx.wsReader()

		// Inform a wait group that this connection has started
		ctx.wg.Done()

		select {
		case <-ctx.Con.Done():

			// The connection has been closed, try to reconnect
			log.Println("Connection closed. Reconnecting...")
			time.Sleep(time.Duration(ctx.SecondPerAttempt) * time.Second) // Wait before reconnecting
		case <-ctx.CloseWS:
			// Close requested, gracefully close the WebSocket
			if err := ctx.C.Close(websocket.StatusNormalClosure, "Client going offline"); err != nil {
				log.Printf("Failed to close WebSocket: %v", err)
			}

		}
	}
}

// NewConnection starts a new Websocket connection
func NewConnection(URI string) *Ws {

	var ctx Ws
	ctx.URI = URI
	ctx.wg.Add(1)
	go ctx.startWsConnection()
	ctx.wg.Wait()
	return &ctx

}

// QueryByID queries HSVR ELO bot by ID and type
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

	// Executes queries
	ctx.WriteMessage(newMessage)

}

// Subscribe sends a subscribe request to HSVR ELO bot events.
func (ctx *Ws) Subscribe(Type []string) {

	subscribe := WsMessage{
		Type: ctx.subscribeType(),
		Data: ctx.Subscriptions.All(),
		Pid:  uuid.NewString(),
	}

	// Sends a subscribe request.
	ctx.WriteMessage(subscribe)
}

// WriteMessage writes a message to the websocket connection.
func (ctx *Ws) WriteMessage(v any) error {
	ctx.writeWg.Lock()
	err := wsjson.Write(ctx.Con, ctx.C, v)
	if err != nil {
		fmt.Println(err.Error())
		ctx.writeWg.Unlock()
		return err
	}
	ctx.writeWg.Unlock()

	return err
}

// wsReader reads a websocket connection until it is closed.
func (ctx *Ws) wsReader() {
	for {

		// Read from the WebSocket connection
		_, wsMsgByte, err := ctx.C.Reader(ctx.Con)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// Start a Goroutine to handle the message concurrently
		//go ctx.handleWebSocketMessage(wsMsgByte)
		message, err := ctx.decodeWsMessage(wsMsgByte)
		fmt.Println(message)

		if err != nil {
			// ...
		}

		if ctx.last == message {
			continue
		}

		ctx.last = message
		go ctx.handleWsMessage(message)
	}
}

func (ctx *Ws) decodeWsMessage(wsMsgByte io.Reader) (string, error) {
	buffer := make([]byte, 1024)
	var message string

	// Read a chunk of data from the WebSocket
	for {
		count, err := wsMsgByte.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return message, err
			}
			fmt.Println(err.Error())
			return "", err
		}

		// Read a chunk of data from the WebSocket
		message = string(buffer[:count])

	}
}

func (ctx *Ws) handleWsMessage(message string) {

	// Ping?
	if isPingMessage(message) {
		// Pong!

		fmt.Println("Pong!")
		err := ctx.sendPongMessage()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		// ...

		var jsonMsg WsMessage
		err := json.Unmarshal([]byte(message), &jsonMsg)
		if err != nil {
			return
		}
		ctx.assertWsMessageType(message)

		found := false
		for _, x := range ctx.TypesFound {
			if x == jsonMsg.Type {
				found = true
			}
		}

		if !found {
			ctx.TypesFound = append(ctx.TypesFound, jsonMsg.Type)
			fmt.Print(jsonMsg.Type + "\n" + message + "\n\n")
		}
	}
}

func (ctx *Ws) assertWsMessageType(message string) {

	wsMessage := WsMessage{}

	err := json.Unmarshal([]byte(message), &wsMessage)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	jsonMarshal, err := json.Marshal(wsMessage.Data)
	if err != nil {
		return
	}

	switch wsMessage.Type {
	case "online":

		for _, x := range ctx.BadUnmarshals {
			if x == "online" {
				break
			}
		}

		for _, x := range ctx.GoodUnmarshals {
			if x == "online" {
				break
			}
		}
		online := []OnlineData{}
		err := json.Unmarshal(jsonMarshal, &online)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		ctx.GoodUnmarshals = append(ctx.GoodUnmarshals, "online")
		fmt.Println("online unmarshal good")
	case "kill":

		for _, x := range ctx.BadUnmarshals {
			if x == "kill" {
				break
			}
		}

		for _, x := range ctx.GoodUnmarshals {
			if x == "kill" {
				break
			}
		}
		kill := global.KillData{}
		err := json.Unmarshal(jsonMarshal, &kill)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		ctx.GoodUnmarshals = append(ctx.GoodUnmarshals, "kill")
		fmt.Println("kill unmarshal good")
		handlers.OnKillStream(kill)
	case "tracking":
		tracking := Tracking{}
		err := json.Unmarshal(jsonMarshal, &tracking)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		found := false
		for _, x := range ctx.TrackingTypesFound {
			if x == tracking.TrackingType {
				found = true
			}
		}

		trackData, err := json.Marshal(tracking.TrackingData)
		if err != nil {
			return
		}

		if !found {
			fmt.Println("newTrackingType: " + tracking.TrackingType + "\n" + string(trackData) + "\n\n")
			ctx.TrackingTypesFound = append(ctx.TrackingTypesFound, tracking.TrackingType)
		}
	case "spawn":
		for _, x := range ctx.BadUnmarshals {
			if x == "spawn" {
				break
			}
		}

		for _, x := range ctx.GoodUnmarshals {
			if x == "spawn" {
				break
			}
		}
		spawn := SpawnData{}
		err := json.Unmarshal(jsonMarshal, &spawn)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		ctx.GoodUnmarshals = append(ctx.GoodUnmarshals, "spawn")
		fmt.Println("spawn unmarshal good")
	case "user_login":
		for _, x := range ctx.BadUnmarshals {
			if x == "user_login" {
				break
			}
		}

		for _, x := range ctx.GoodUnmarshals {
			if x == "user_login" {
				break
			}
		}
		login := UserLogEvent{}
		err := json.Unmarshal(jsonMarshal, &login)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		ctx.GoodUnmarshals = append(ctx.GoodUnmarshals, "user_login")
		fmt.Println("user_login unmarshal good")
	case "user_logout":
		for _, x := range ctx.BadUnmarshals {
			if x == "user_logout" {
				break
			}
		}

		for _, x := range ctx.GoodUnmarshals {
			if x == "user_logout" {
				break
			}
		}
		logout := UserLogEvent{}
		err := json.Unmarshal(jsonMarshal, &logout)
		if err != nil {
			fmt.Println(err.Error())
			ctx.BadUnmarshals = append(ctx.BadUnmarshals, "user_logout")
			return
		}
		ctx.GoodUnmarshals = append(ctx.GoodUnmarshals, "user_logout")
		fmt.Println("user_logout unmarshal good")
	}

}

// isPingMessage checks if a given message is a ping message.
func isPingMessage(message string) bool {
	// Check if the message is a ping message (you can define your criteria)

	return strings.Contains(message, "ping")
}

// sendPongMessage sends a pong message in response to a ping message.
func (ctx *Ws) sendPongMessage() error {
	pong := PongMessage{
		Type: "pong",
		PID:  uuid.NewString(),
	}

	// Pong!
	err := ctx.WriteMessage(pong)
	return err
}
