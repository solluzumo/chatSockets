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

type ChatRepoPostgres struct {
	db *gorm.DB
}

func NewChatRepoPostgres(db *gorm.DB) *ChatRepoPostgres {
	return &ChatRepoPostgres{
		db: db,
	}
}

func (cr *ChatRepoPostgres) Count(ctx context.Context) int64 {
	var count int64
	cr.db.WithContext(ctx).Model(&models.Chat{}).Where("1=1").Count(&count)
	return count
}

func (cr *ChatRepoPostgres) CreateChat(ctx context.Context, data *domain.ChatDomain) error {
	var err error

	data.CreatedAt = time.Now()

	chatModel := models.ToModelChat(data)

	if err = cr.db.WithContext(ctx).Create(chatModel).Error; err != nil {
		if IsUniqueViolation(err) {
			err = domain.ErrChatAlreadyExists
		}
		return fmt.Errorf("не удалось создать чат: %w", err)
	}

	data.ID = chatModel.ID

	return nil
}

func (cr *ChatRepoPostgres) GetChatWithUsers(ctx context.Context, chatID int) (*domain.ChatDomain, error) {
	// Сначала достаём чат
	var chat models.Chat
	if err := cr.db.WithContext(ctx).First(&chat, chatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("чат с ID %d не найден", chatID)
		}
		return nil, fmt.Errorf("ошибка при получении чата: %w", err)
	}

	// Потом всех участников
	var userLinks []models.UserChatLink
	if err := cr.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Find(&userLinks).Error; err != nil {
		return nil, fmt.Errorf("ошибка при получении участников чата: %w", err)
	}

	// Маппим в доменную модель
	chatDomain := &domain.ChatDomain{
		ID:         chat.ID,
		Title:      chat.Title,
		CreatedAt:  chat.CreatedAt,
		ChatStatus: chat.ChatStatus,
		Users:      make([]*domain.UserChatLinkDomain, len(userLinks)),
	}

	for i, u := range userLinks {
		chatDomain.Users[i] = &domain.UserChatLinkDomain{
			UserID:      u.UserID,
			ChatID:      u.ChatID,
			UserBlocked: u.UserBlocked,
			UserRole:    u.UserRole,
			CreatedAt:   u.CreatedAt,
		}
	}

	return chatDomain, nil
}

func (cr *ChatRepoPostgres) ChatExists(ctx context.Context, chatDomain *domain.ChatDomain) (bool, error) {
	var count int64

	chatModel := models.ToModelChat(chatDomain)

	cr.db.WithContext(ctx).Model(chatModel).Where(chatModel).Count(&count)

	return count > 0, nil
}

func (cr *ChatRepoPostgres) DeleteChat(ctx context.Context, chatID int) error {

	result := cr.db.WithContext(ctx).Where("id = ?", chatID).Delete(&models.Chat{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrChatNotFound
	}
	return nil
}

// ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ

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
