package app

import (
	"chatsockets/internal/events"
	httpHandlers "chatsockets/internal/interfaces/handlers"
	"chatsockets/internal/repo/postgres"
	"chatsockets/internal/services"
	websockets "chatsockets/internal/ws"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppHandlers struct {
	MessageHandler *httpHandlers.MessageHandler
	ChatHandler    *httpHandlers.ChatHandler
	LinkHandler    *httpHandlers.LinkHandler
}

type AppServices struct {
	PermissionService *services.PermissionService
	MessageHub        *websockets.MessageHub
	MessageService    *services.MessageService
	ChatService       *services.ChatService
	LinkService       *services.LinkService
}

type AppRepos struct {
	ChatRepo    *postgres.ChatRepoPostgres
	MessageRepo *postgres.MessageRepoPostgres
	LinkRepo    *postgres.LinkPostgresRepo
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
	MsgChan chan events.MessageCreatedEvent) *AppInstance {

	eventBus := events.NewEventBus()

	appRepos := &AppRepos{
		ChatRepo:    postgres.NewChatRepoPostgres(db),
		MessageRepo: postgres.NewMessageRepoPostgres(db),
		LinkRepo:    postgres.NewLinkPostgresRepo(db),
	}

	permissionService := services.NewPermissionService(appRepos.LinkRepo)

	appServices := &AppServices{
		PermissionService: permissionService,
		MessageHub:        websockets.NewMessageHub(WG, appLogger, eventBus),
		MessageService:    services.NewMessageService(appRepos.MessageRepo, appRepos.ChatRepo, appRepos.LinkRepo, eventBus, permissionService),
		ChatService:       services.NewChatService(appRepos.MessageRepo, appRepos.ChatRepo, appRepos.LinkRepo, permissionService),
		LinkService:       services.NewLinkService(appRepos.LinkRepo, permissionService),
	}

	appHandlers := &AppHandlers{
		MessageHandler: httpHandlers.NewMessageHandler(upgrader, appServices.MessageService, appLogger, appServices.MessageHub),
		ChatHandler:    httpHandlers.NewChatAPIHTTP(appServices.ChatService, appLogger),
		LinkHandler:    httpHandlers.NewLinkAPIHTTP(appServices.LinkService, appLogger),
	}

	return &AppInstance{
		Handlers: appHandlers,
		Services: appServices,
		Repos:    appRepos,
	}
}
