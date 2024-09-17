package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/diwasrimal/echo-backend/db"
	"github.com/diwasrimal/echo-backend/models"
	"github.com/diwasrimal/echo-backend/types"
	"github.com/gorilla/websocket"
)

// Stores the connection of each client
var clientsMu sync.RWMutex
var clients = make(map[uint64]*websocket.Conn)

var wsUp = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type payloadMsgType string

const (
	chatMsgSend    payloadMsgType = "chatMsgSend"
	chatMsgReceive payloadMsgType = "chatMsgReceive"
)

type wsPayload struct {
	MsgType payloadMsgType `json:"msgType"`
	MsgData any            `json:"msgData"`
}

func WSHandleFunc(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint64)
	log.Printf("Hit WSHandleFunc() with userId: %v\n", userId)
	conn, err := wsUp.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to ws: %v\n", err)
		return
	}

	clientsMu.Lock()
	clients[userId] = conn
	clientsMu.Unlock()

	// Remove the connection from map and close
	defer func() {
		clientsMu.Lock()
		delete(clients, userId)
		clientsMu.Unlock()
		conn.Close()
	}()

	for true {
		var payload wsPayload
		err := conn.ReadJSON(&payload)
		if err != nil {
			log.Printf("%T reading ws json data: %v\n", err, err)
			break
		}

		log.Printf("ws payload: %+v\n", payload)

		switch payload.MsgType {
		case chatMsgSend:
			err = handleChatMsgSend(userId, payload.MsgData.(map[string]any))
		}
		if err != nil {
			log.Printf("Error handling %v: %v\n", payload.MsgType, err)
		}
	}
}

func handleChatMsgSend(senderId uint64, data types.Json) error {
	rid, ridOk := data["receiverId"].(float64)
	text, textOk := data["text"].(string)
	ts, tsOk := data["timestamp"].(string)
	if !ridOk || !textOk || !tsOk {
		return errors.New("Invalid/Missing data fields")
	}
	if len(text) == 0 {
		return nil
	}

	receiverId := uint64(rid)
	timestamp, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return err
	}

	// Store the message in database
	msgId, err := db.RecordMessage(senderId, receiverId, text, timestamp)
	if err != nil {
		return fmt.Errorf("Can't record msg of %v in db: %v", senderId, err)
	}

	// Update the last conversation time for two users
	if err := db.UpdateOrCreateConversation(senderId, receiverId, timestamp); err != nil {
		return fmt.Errorf("Can't update last conv for (%v,%v): %v\n", senderId, receiverId, err)
	}

	// Broadcast the chat message to both sender and receiver
	payload := wsPayload{
		MsgType: chatMsgReceive,
		MsgData: models.Message{
			Id:         msgId,
			SenderId:   senderId,
			ReceiverId: receiverId,
			Text:       text,
			Timestamp:  timestamp,
		},
	}

	clientsMu.RLock()
	senderConn, ok := clients[senderId]
	clientsMu.RUnlock()
	if ok {
		err := senderConn.WriteJSON(payload)
		if err != nil {
			log.Printf("Error writing json to msg sender: %v\n", err)
		}
	}

	clientsMu.RLock()
	receiverConn, ok := clients[receiverId]
	clientsMu.RUnlock()
	if ok {
		err := receiverConn.WriteJSON(payload)
		if err != nil {
			log.Printf("Error writing json to msg receiver: %v\n", err)
		}
	}

	return nil
}
