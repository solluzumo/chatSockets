package services

import (
	"chatsockets/internal/dto"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

func Worker(id int, ctx context.Context, wg *sync.WaitGroup, msgChannel chan []byte, clients map[*websocket.Conn]bool, rwMut *sync.RWMutex) {
	defer wg.Done()
	log.Printf("Воркер %d запущен\n", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Воркер %d завершается по контексту", id)
			return

		case msg, ok := <-msgChannel:
			if !ok {
				log.Printf("Воркер %d: канал закрыт, выходим", id)
				return
			}

			log.Printf("Воркер %d обработал сообщение", id)

			var messageSendRequest dto.MessageSendRequest
			if err := json.Unmarshal(msg, &messageSendRequest); err != nil {
				log.Printf("ошибка анмаршелинга: %v", err)
				continue
			}

			// Блокируем чтение мапы, чтобы никто не менял её в этот момент
			rwMut.RLock()
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Printf("ошибка отправки клиенту: %v", err)
				}
			}
			rwMut.RUnlock()
		}
	}
}
