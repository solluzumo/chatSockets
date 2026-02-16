package domain

type UserAction string

const (
	//Сообщения
	SendMessage      UserAction = "отправить сообщение"
	DeleteMessage    UserAction = "удалить чужое сообщение"
	UpdateMessage    UserAction = "изменить чужое сообщение"
	DeleteOwnMessage UserAction = "удалить своё сообщение"
	UpdateOwnMessage UserAction = "изменить своё сообщение"

	//Чат
	SubscribeToChat UserAction = "подписать на чат" // на будущее когда чат будет закрытым
	CheckChat       UserAction = "посмотреть чат"
	DeleteChat      UserAction = "удалить чат"
	UpdateChat      UserAction = "изменить чат"

	//Другие пользователи в чате
	UpdateRole UserAction = "изменить роль"
	BanUser    UserAction = "заблокировать пользователя"
	AddUser    UserAction = "добавить пользователя"
)
