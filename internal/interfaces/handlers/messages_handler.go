package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWaitS = 10 * time.Second
	pingWaitS = (pongWaitS * 9) / 10
	writeWait = 10 * time.Second
)

type MessageHandler struct {
	upgrader   websocket.Upgrader
	clients    map[*websocket.Conn]bool
	msgChannel chan []byte
	rwmut      *sync.RWMutex
	wg         *sync.WaitGroup
}

func NewMessageHandler(upgrader websocket.Upgrader, clients map[*websocket.Conn]bool, msgChan chan []byte, rwmut *sync.RWMutex, wg *sync.WaitGroup) *MessageHandler {
	return &MessageHandler{
		upgrader:   upgrader,
		clients:    clients,
		msgChannel: msgChan,
		rwmut:      rwmut,
		wg:         wg,
	}
}

func (mh *MessageHandler) deleteClient(rw *sync.RWMutex, client *websocket.Conn) {
	rw.RLock()
	defer rw.RUnlock()

	delete(mh.clients, client)
}

func (mh *MessageHandler) GetAndPostMessage(w http.ResponseWriter, r *http.Request) {

	ws, err := mh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка при создании подключения%v", err)
		return
	}

	mh.clients[ws] = true

	go mh.writePump(ws)
	go mh.readPump(ws)
}

func (mh *MessageHandler) readPump(conn *websocket.Conn) {
	defer conn.Close()

	//ставим дедлайн для проверки подключения
	conn.SetReadDeadline(time.Now().Add(pongWaitS))

	//обновляем счётчик после получения понга
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWaitS))
		return nil
	})

	for {
		// Читаем сообщение из сокета
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}
		mh.msgChannel <- message
	}

}

func (mh *MessageHandler) writePump(conn *websocket.Conn) {
	ticker := time.NewTicker(pingWaitS)

	defer func() {
		ticker.Stop()
		conn.Close()
		mh.deleteClient(mh.rwmut, conn)
	}()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Printf("write control error: %v", err)
				return
			}

		}
	}
}
