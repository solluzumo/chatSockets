package postgres

import (
	"chatsockets/internal/domain"
	"chatsockets/internal/models"
	"chatsockets/internal/repo"
	"context"
	"errors"
	"fmt"

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

func (cr *ChatRepoPostgres) CreateChat(ctx context.Context, data *domain.ChatDomain) (*domain.ChatDomain, error) {
	var err error

	chatModel := models.DomainToModelChat(data)

	if err = cr.db.WithContext(ctx).Create(chatModel).Error; err != nil {
		if IsUniqueViolation(err) {
			err = domain.ErrChatAlreadyExists
		}
		return data, fmt.Errorf("не удалось создать чат: %w", err)
	}

	data.ID = chatModel.ID

	return data, nil
}

func (cr *ChatRepoPostgres) FindChatById(ctx context.Context, data *domain.ChatDomain) error {
	chatModel := &models.Chat{}

	err := cr.db.WithContext(ctx).First(&chatModel, data.ID).Error
	if err != nil {
		return fmt.Errorf("не удалось найти чат:%w", err)
	}

	data.Title = chatModel.Title

	data.ID = chatModel.ID

	return nil
}

func (cr *ChatRepoPostgres) ChatExists(ctx context.Context, param repo.FilterParam) (bool, error) {
	var count int64

	chatModel := &models.Chat{}

	if isFieldAllowed(param.Field) != nil {
		return false, fmt.Errorf("не удалось проверить существует ли чат: %w", domain.ErrFieldIsNotAllowed)
	}

	cr.db.WithContext(ctx).Model(chatModel).Where(param.Field+" = ?", param.Value).Count(&count)

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
