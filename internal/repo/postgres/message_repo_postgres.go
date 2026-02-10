package postgres

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MessageRepoPostgres struct {
	db *gorm.DB
}

func NewMessageRepoPostgres(db *gorm.DB) *MessageRepoPostgres {
	return &MessageRepoPostgres{
		db: db,
	}
}

func (mr *MessageRepoPostgres) Count(ctx context.Context) int64 {
	var count int64
	mr.db.WithContext(ctx).Model(&models.Message{}).Where("1=1").Count(&count)
	return count
}

func (mr *MessageRepoPostgres) CreateMessage(ctx context.Context, data *domain.MessageDomain) (*domain.MessageDomain, error) {
	var err error

	messageModel := models.DomainToModelMessage(data)

	messageModel.CreatedAt = time.Now()

	if err = mr.db.WithContext(ctx).Create(messageModel).Error; err != nil {
		if pgErr, ok := IsForeignKeyViolation(err); ok {
			return nil, MapFKConstraint(pgErr.ConstraintName)
		}
		return nil, fmt.Errorf("ошибка создания сообщения: %w", err)
	}

	data.ID = messageModel.ChatID

	return data, nil
}

func (mr *MessageRepoPostgres) GetMessagesByChatWithLimit(ctx context.Context, chatID int, limit int) ([]*domain.MessageDomain, error) {
	var messageModels []*models.Message

	if err := mr.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&messageModels).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("не удалось получить сообщения из чата: %w", domain.ErrChatNotFound)
		}
		return nil, fmt.Errorf("не удалось получить сообщения из чата: %w", err)
	}
	messageDomains := models.ModelSliceToDomainSlice(messageModels)

	return messageDomains, nil
}

func (mr *MessageRepoPostgres) DeleteMessages(ctx context.Context, chatID int) error {
	var err error
	if err = mr.db.WithContext(ctx).Where("chat_id = ?", chatID).Delete(&models.Message{}).Error; err != nil {
		if IsUniqueViolation(err) {
			err = domain.ErrMessageNotFound
		}
		return fmt.Errorf("не удалось удалить сообщение: %w", err)
	}
	return nil
}
