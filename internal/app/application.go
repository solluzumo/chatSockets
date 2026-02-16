package app

import (
	"chatsockets/internal/events"
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	MsgChannel chan events.MessageCreatedEvent
	RWMutext   *sync.RWMutex
	WG         *sync.WaitGroup
	Instance   *AppInstance
	Cfg        *Config
}

func NewApplication(
	upgrader websocket.Upgrader,
	clients map[*websocket.Conn]bool,
	db *gorm.DB,
	appLogger *zap.Logger) *Application {

	var rwMutext sync.RWMutex
	var wg sync.WaitGroup

	cfg := NewConfig()
	if cfg == nil {
		return nil
	}
	msgChannel := make(chan events.MessageCreatedEvent, cfg.MessageChannelCap)

	return &Application{
		MsgChannel: msgChannel,
		RWMutext:   &rwMutext,
		WG:         &wg,
		Instance:   NewAppInstance(upgrader, clients, db, appLogger, &rwMutext, &wg, msgChannel),
		Cfg:        cfg,
	}
}

func (a *Application) Start(ctx *context.Context, cfg *Config) {
	//ВОРКЕРЫ
	for i := range a.Cfg.MessageWorkersCap {
		a.WG.Add(1)
		go a.Instance.Services.MessageHub.MessageWorker(i, *ctx)
	}
}
