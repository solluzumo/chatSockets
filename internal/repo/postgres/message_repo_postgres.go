package postgres

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/models"
	"context"

	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MessageRepoPostgres struct {
	db       *gorm.DB
	dbLogger *zap.Logger
}

func NewMessageRepoPostgres(db *gorm.DB, appLogger *zap.Logger) *MessageRepoPostgres {
	dbLogger := appLogger.Named("message_db")
	return &MessageRepoPostgres{
		db:       db,
		dbLogger: dbLogger,
	}
}

func (mr *MessageRepoPostgres) Count(ctx context.Context) int64 {
	var count int64
	mr.db.WithContext(ctx).Model(&models.Message{}).Where("1=1").Count(&count)
	return count
}

func (mr *MessageRepoPostgres) CreateMessage(ctx context.Context, data *domain.MessageDomain) (*domain.MessageDomain, error) {
	messageModel := &models.Message{}

	messageModel.ChatID = data.ChatID
	messageModel.Text = data.Text
	messageModel.CreatedAt = time.Now()

	if err := mr.db.WithContext(ctx).Create(messageModel).Error; err != nil {
		return data, err
	}

	mr.dbLogger.Info("Сообщение создано", zap.Int("message id", messageModel.ID))

	data.ID = messageModel.ChatID

	return data, nil
}

func (mr *MessageRepoPostgres) GetMessagesByChatWithLimit(ctx context.Context, chatID int, limit int) []*domain.MessageDomain {
	var messageModels []*models.Message

	var messageDomains []*domain.MessageDomain

	if err := mr.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&messageModels).Error; err != nil {
		return nil
	}

	for _, el := range messageModels {
		messageDomains = append(messageDomains, &domain.MessageDomain{
			ID:        el.ID,
			ChatID:    el.ChatID,
			Text:      el.Text,
			CreatedAt: el.CreatedAt,
		})
	}

	return messageDomains
}

func (mr *MessageRepoPostgres) DeleteMessages(ctx context.Context, chatID int) error {
	return mr.db.WithContext(ctx).Where("chat_id = ?", chatID).Delete(&models.Message{}).Error
}
