package database

type User struct {
	UserId   int64  `json:"userId"`
	Name     string `json:"name"`
	PhotoUrl string `json:"photoUrl"`
}

type Group struct {
	ConversationId int64  `json:"conversationId"`
	FounderId      int64  `json:"founderId"`
	Name           string `json:"name"`
	PhotoUrl       string `json:"photoUrl"`
	CreatedAt      string `json:"createdAt"`
}

type UserConversation struct {
	ConversationId int64 `json:"conversationId"`
	UserId         int64 `json:"userId"`
}

type Message struct {
	MessageId      int64     `json:"messageId"`
	Text           string    `json:"text"`
	SenderId       int64     `json:"senderId"`
	ConversationId int64     `json:"conversationId"`
	CreatedAt      string    `json:"createdAt"`
	IsForwarded    bool      `json:"isForwarded"`
	CommentTo      *int64    `json:"commentTo"`
	AttachmentUrl  string    `json:"attachmentUrl"`
	MediaType      MediaType `json:"mediaType"`
}

type Receipt struct {
	MessageId  int64   `json:"messageId"`
	UserId     int64   `json:"userId"`
	Status     Status  `json:"status"`
	SentAt     string  `json:"sentAt"`
	ReceivedAt *string `json:"receivedAt"`
	ReadAt     *string `json:"readAt"`
}

type Reaction struct {
	Emoji     Emoji `json:"emoji"`
	MessageId int64 `json:"messageId"`
	SenderId  int64 `json:"senderId"`
}
