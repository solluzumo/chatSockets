package main

import (
	"chatsockets/internal/app"
	"chatsockets/internal/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func main() {
	//ДЛЯ ВЕБСОКЕТОВ
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
	}
	clients := make(map[*websocket.Conn]bool)

	//БАЗОВЫЙ ЛОГГЕР
	logger := app.NewZapLogger()
	defer logger.Sync()

	//ГЛОБАЛЬНЫЙ ЛОГГЕР ПРИЛОЖЕНИЯ
	appLogger := logger.Named("app")

	//БАЗА ДАННЫХ
	db, err := app.InitDb()
	if err != nil {
		appLogger.Error("не удалось подключиться к базе данных: ", zap.Error(err))
		panic(err)
	}
	appLogger.Info("сервер подключен к базе данных", zap.Time("started", time.Now()))

	//КОНТЕКСТ ДЛЯ GRACEFULL SHUTDOWN
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	//ОСНОВНАЯ СУЩНОСТЬ ПРИЛОЖЕНИЯ DI
	application := app.NewApplication(upgrader, clients, db, appLogger)
	if application == nil {
		appLogger.Error("не удалось инициализировать Application")
		panic(fmt.Errorf("не удалось инициализировать Application"))
	}

	//Запускаем фоновые задачи(воркеров)
	application.Start(&ctx, application.Cfg)

	//ПОЛУЧАЕМ PUBLIC KEY JWT
	publicKey, err := app.LoadRSAPublicKey(application.Cfg.JWTPublicKeyPath)
	if err != nil {
		appLogger.Error("не удалось загрузить public key: ", zap.Error(err))
		panic(err)
	}

	//ИНИЦИАЛИЗИРУЕМ AUTH MIDDLEWARE
	authMiddleware := middleware.NewJWTMiddleware(publicKey, application.Cfg.JWTIssuer, appLogger)

	//ИНИЦИАЛИЗИРУЕМ И НАСТРАИВАЕМ РОУТЕР
	router := chi.NewRouter()

	app.RegisterRoutes(router, application.Instance, appLogger, authMiddleware)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	fmt.Println("WebSocket server started on :8080")
	//ЗАПУСКАЕМ СЕРВЕР
	go func() {
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

	close(application.MsgChannel)
	application.WG.Wait()

	fmt.Println("Server exited properly")
}
