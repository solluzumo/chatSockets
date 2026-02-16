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

type LinkPostgresRepo struct {
	db *gorm.DB
}

func NewLinkPostgresRepo(db *gorm.DB) *LinkPostgresRepo {
	return &LinkPostgresRepo{
		db: db,
	}
}

func (cr *LinkPostgresRepo) UpdateLink(ctx context.Context, linkDomain *domain.UserChatLinkDomain) error {
	linkModel := models.ToUserChatLink(linkDomain)
	result := cr.db.Model(&models.UserChatLink{}).Where("user_id = ? AND chat_id = ?", linkModel.UserID, linkModel.ChatID).Updates(linkModel)
	if result.Error != nil {
		return fmt.Errorf("не удалось изменить роль: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNoChanges
	}
	return nil
}

func (cr *LinkPostgresRepo) CreateLink(ctx context.Context, linkDomain *domain.UserChatLinkDomain) error {
	var err error

	linkModel := models.ToUserChatLink(linkDomain)

	linkModel.CreatedAt = time.Now()

	if err = cr.db.WithContext(ctx).Create(linkModel).Error; err != nil {
		if IsUniqueViolation(err) {
			err = domain.ErrUserAlreadySubscribed
		}
		return fmt.Errorf("не удалось подписать пользователя к чату: %w", err)
	}

	return nil
}

func (cr *LinkPostgresRepo) GetLink(ctx context.Context, linkDomain *domain.UserChatLinkDomain) error {
	var err error
	linkModel := models.ToUserChatLink(linkDomain)

	err = cr.db.WithContext(ctx).Where(&linkModel).First(&linkModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = domain.ErrUserIsNotConnectedToChat
		}
		return fmt.Errorf("не удалось найти связь пользователь-чат: %w", err)
	}

	//Заполняем домен новыми данными(роли, блокировки и тд)
	linkModel.UpdateUserChatLinkDomain(linkDomain)

	return nil
}

func (cr *LinkPostgresRepo) LinkExists(ctx context.Context, linkDomain *domain.UserChatLinkDomain) (bool, error) {
	var count int64

	linkModel := models.ToUserChatLink(linkDomain)

	if err := cr.db.WithContext(ctx).Model(linkModel).Where(linkModel).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
