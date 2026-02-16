package domain

import "fmt"

type UserRole string

const (
	AdminRole  UserRole = "Администратор"
	GuestRole  UserRole = "Гость"
	MemberRole UserRole = "Участник"
)

func ParseUserRole(role string) (UserRole, error) {
	switch role {
	case string(AdminRole):
		return AdminRole, nil
	case string(MemberRole):
		return MemberRole, nil
	case string(GuestRole):
		return GuestRole, nil
	default:
		return "", fmt.Errorf("неизвестная роль: %s", role)
	}
}
