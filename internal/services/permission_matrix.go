package services

import "chatsockets/internal/domain"

var roleMatrix = map[domain.UserRole]map[domain.UserAction]bool{
	domain.AdminRole: {
		domain.SendMessage:      true,
		domain.DeleteMessage:    true,
		domain.UpdateMessage:    true,
		domain.DeleteOwnMessage: true,
		domain.UpdateOwnMessage: true,

		domain.SubscribeToChat: true,
		domain.CheckChat:       true,
		domain.DeleteChat:      true,
		domain.UpdateChat:      true,

		domain.UpdateRole: true,
		domain.BanUser:    true,
		domain.AddUser:    true,
	},
	domain.GuestRole: {
		domain.SendMessage:      false,
		domain.DeleteMessage:    false,
		domain.UpdateMessage:    false,
		domain.DeleteOwnMessage: false,
		domain.UpdateOwnMessage: false,

		domain.SubscribeToChat: true,
		domain.CheckChat:       true,
		domain.DeleteChat:      false,
		domain.UpdateChat:      false,

		domain.UpdateRole: false,
		domain.BanUser:    false,
		domain.AddUser:    false,
	},
	domain.MemberRole: {
		domain.SendMessage:      true,
		domain.DeleteMessage:    false,
		domain.UpdateMessage:    false,
		domain.DeleteOwnMessage: true,
		domain.UpdateOwnMessage: true,

		domain.SubscribeToChat: true,
		domain.CheckChat:       true,
		domain.DeleteChat:      false,
		domain.UpdateChat:      false,

		domain.UpdateRole: false,
		domain.BanUser:    false,
		domain.AddUser:    false,
	},
}

var chatMatrix = map[domain.ChatStatus]map[domain.UserAction]bool{
	domain.PrivateChat: {
		domain.SendMessage:      true,
		domain.DeleteMessage:    true,
		domain.UpdateMessage:    true,
		domain.DeleteOwnMessage: true,
		domain.UpdateOwnMessage: true,

		domain.SubscribeToChat: false,
		domain.CheckChat:       true,
		domain.DeleteChat:      true,
		domain.UpdateChat:      true,

		domain.UpdateRole: true,
		domain.BanUser:    true,
		domain.AddUser:    true,
	},
	domain.PublicChat: {
		domain.SendMessage:      true,
		domain.DeleteMessage:    true,
		domain.UpdateMessage:    true,
		domain.DeleteOwnMessage: true,
		domain.UpdateOwnMessage: true,

		domain.SubscribeToChat: true,
		domain.CheckChat:       true,
		domain.DeleteChat:      true,
		domain.UpdateChat:      true,

		domain.UpdateRole: true,
		domain.BanUser:    true,
		domain.AddUser:    true,
	},
	domain.ChannelChan: {
		domain.SendMessage:      true,
		domain.DeleteOwnMessage: true,
		domain.UpdateOwnMessage: true,

		domain.SubscribeToChat: true,
		domain.CheckChat:       true,
		domain.DeleteChat:      true,
		domain.UpdateChat:      true,

		domain.UpdateRole: true,
		domain.BanUser:    true,
		domain.AddUser:    true,
	},
}
