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
	MessageId      int64  `json:"messageId"`
	Text           string `json:"text"`
	AttachmentId   *int64 `json:"attachmentId"`
	SenderId       int64  `json:"senderId"`
	ConversationId int64  `json:"conversationId"`
	CreatedAt      string `json:"createdAt"`
	IsForwarded    bool   `json:"isForwarded"`
	CommentTo      *int64 `json:"commentTo"`
}

type Attachment struct {
	AttachmentId int64  `json:"attachmentId"`
	Url          string `json:"url"`
	MediaType    string `json:"mediaType"`
}

type Status struct {
	MessageId int64 `json:"messageId"`
	UserId    int64 `json:"userId"`
	Info      Info  `json:"info"`
}

type Reaction struct {
	EmojiCode string `json:"emojiCode"`
	MessageId int64  `json:"messageId"`
	SenderId  int64  `json:"senderId"`
}
