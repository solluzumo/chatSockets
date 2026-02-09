package postgres

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/models"
	"chatsockets/internal/repo"
	"context"
	"errors"

	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChatRepoPostgres struct {
	db       *gorm.DB
	dbLogger *zap.Logger
}

func NewChatRepoPostgres(db *gorm.DB, appLogger *zap.Logger) *ChatRepoPostgres {
	dbLogger := appLogger.Named("chat_db")
	return &ChatRepoPostgres{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (cr *ChatRepoPostgres) Count(ctx context.Context) int64 {
	var count int64
	cr.db.WithContext(ctx).Model(&models.Chat{}).Where("1=1").Count(&count)
	return count
}

func (cr *ChatRepoPostgres) CreateChat(ctx context.Context, data *domain.ChatDomain) (*domain.ChatDomain, error) {
	chatModel := &models.Chat{
		Title:     data.Title,
		CreatedAt: time.Now(),
	}

	if err := cr.db.WithContext(ctx).Create(chatModel).Error; err != nil {
		return data, err
	}

	data.ID = chatModel.ID

	return data, nil
}

func (cr *ChatRepoPostgres) FindChatById(ctx context.Context, data *domain.ChatDomain) error {
	chatModel := &models.Chat{}

	err := cr.db.WithContext(ctx).First(&chatModel, data.ID).Error
	if err != nil {
		return err
	}

	data.Title = chatModel.Title

	data.ID = chatModel.ID

	return nil
}

func (cr *ChatRepoPostgres) ChatExists(ctx context.Context, param repo.FilterParam) (bool, error) {
	var count int64

	chatModel := &models.Chat{}

	if isFieldAllowed(param.Field) != nil {
		return false, domain.ErrFieldIsNotAllowed
	}

	cr.db.WithContext(ctx).Model(chatModel).Where(param.Field+" = ?", param.Value).Count(&count)

	return count > 0, nil
}

func (cr *ChatRepoPostgres) DeleteChat(ctx context.Context, chatID int) error {
	return cr.db.WithContext(ctx).Where("id = ?", chatID).Delete(&models.Chat{}).Error
}

func (cr *ChatRepoPostgres) IsUserConnectedToChat(ctx context.Context, chatID int, userID int) bool {
	return true
}

func isFieldAllowed(field string) error {
	allowedFields := map[string]bool{
		"title":      true,
		"created_at": true,
		"id":         true,
	}

	if !allowedFields[field] {
		return errors.New("invalid field")
	}

	return nil
}
