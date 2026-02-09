package app

import (
	"chatsockets/internal/domain"
	httpHandlers "chatsockets/internal/interfaces/handlers"
	"chatsockets/internal/repo/postgres"
	"chatsockets/internal/services"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppHandlers struct {
	MessageHandler *httpHandlers.MessageHandler
	ChatHandler    *httpHandlers.ChatHandler
}

type AppServices struct {
	MessageHub  *services.MessageHub
	ChatService *services.ChatService
}

type AppRepos struct {
	ChatRepo    *postgres.ChatRepoPostgres
	MessageRepo *postgres.MessageRepoPostgres
}

func NewAppRepos() *AppRepos {
	return &AppRepos{}
}

type AppInstance struct {
	Handlers *AppHandlers
	Services *AppServices
	Repos    *AppRepos
}

func NewAppInstance(
	upgrader websocket.Upgrader,
	clients map[*websocket.Conn]bool,
	db *gorm.DB,
	appLogger *zap.Logger,
	RWMutex *sync.RWMutex,
	WG *sync.WaitGroup,
	MsgChan chan domain.MessageTask) *AppInstance {

	AppRepos := &AppRepos{
		ChatRepo:    postgres.NewChatRepoPostgres(db, appLogger),
		MessageRepo: postgres.NewMessageRepoPostgres(db, appLogger),
	}

	appServices := &AppServices{
		MessageHub:  services.NewMessageHub(RWMutex, clients, WG, appLogger, AppRepos.MessageRepo, AppRepos.ChatRepo),
		ChatService: services.NewChatService(AppRepos.MessageRepo, AppRepos.ChatRepo, appLogger),
	}

	appHandlers := &AppHandlers{
		MessageHandler: httpHandlers.NewMessageHandler(upgrader, MsgChan, appServices.MessageHub, appLogger),
		ChatHandler:    httpHandlers.NewChatAPIHTTP(appServices.ChatService, appLogger),
	}

	return &AppInstance{
		Handlers: appHandlers,
		Services: appServices,
		Repos:    AppRepos,
	}
}
