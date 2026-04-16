package service

import (
	"os"
	"time"
)

type User struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Photo os.File `json:"photo"`
}

type Group struct {
	ConversationId *int64 `json:"conversationId"`
	UserId int64 `json:"userId"`
	Name string `json:"name"`
	Photo os.File `json:"photo"`
	Timestamp time.Time `json:"timestamp"`
}

type Message struct {
	Id int64 `json:"id"`
	Content MessageContent `json:"content"`
	UserId int64 `json:"userId"`
	ConversationId int64 `json:"conversationId"`
	Timestamp time.Time `json:"timestamp"`
	IsForwarded bool `json:"isForwarded"`
	MessageId int64 `json:"messageId"`
}

type MessageContent struct {
	Image os.File
	Text string
}

type Reaction struct {
	Emoji int64 `json:"emoji"`
	MessageId int64 `json:"messageId"`
	UserId int64 `json:"userId"`
}
