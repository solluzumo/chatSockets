package main

import (
	"chatsockets/internal/app"
	"chatsockets/internal/interfaces/handlers"
	"chatsockets/internal/services"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MessageChannelCap = 10
	MessageWorkersCap = 10
)

func main() {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
	}

	var rwMutext sync.RWMutex
	var wg sync.WaitGroup
	clients := make(map[*websocket.Conn]bool)
	msgChannel := make(chan []byte, MessageChannelCap)
	app := app.NewApp(msgChannel, &rwMutext, &wg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for i := range MessageWorkersCap {
		app.WG.Add(1)
		go services.Worker(i, ctx, app.WG, app.MsgChannel, clients, app.RWMutext)
	}

	messageHandler := handlers.NewMessageHandler(upgrader, clients, app.MsgChannel, app.RWMutext, app.WG)

	http.HandleFunc("/ws", messageHandler.GetAndPostMessage)
	fmt.Println("WebSocket server started on :8080")

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		fmt.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ждём Ctrl+C
	<-ctx.Done()
	fmt.Println("\nShutting down...")

	// Завершаем HTTP-сервер
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	close(app.MsgChannel)
	wg.Wait()

	fmt.Println("Server exited properly")
}
