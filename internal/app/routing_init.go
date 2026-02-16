package app

import (
	"chatsockets/internal/middleware"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func RegisterRoutes(r chi.Router, app *AppInstance, appLogger *zap.Logger, authMiddleware *middleware.JWTMiddleware) {
	r.Use(middleware.CORSMiddleware)

	r.With(authMiddleware.Handle).
		Group(func(r chi.Router) {
			r.Use(middleware.LoggingMiddleWare(appLogger))
			r.Route("/chats", func(r chi.Router) {
				r.Post("/", app.Handlers.ChatHandler.CreateChat)
				r.Delete("/{id}", app.Handlers.ChatHandler.DeleteChat)
				r.Get("/{id}", app.Handlers.ChatHandler.GetChat)
				r.Post("/{id}/subscribe", app.Handlers.LinkHandler.SubscribeToChat)
				r.Post("/{id}/update-role", app.Handlers.LinkHandler.UpdateUserRole)
				r.Post("/{id}/add", app.Handlers.LinkHandler.AddUserToChat)
			})
		})

	r.With(authMiddleware.HandleWS).
		Get("/chats/{id}/messages/", app.Handlers.MessageHandler.ConnectToChatWS)
}
