package database

import "os"


type User struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Photo *os.File `json:"photo"`
}

type Group struct {
	Conversation *int64 `json:"conversation"`
	Founder int64 `json:"founder"`
	Name string `json:"name"`
	Photo os.File `json:"photo"`
	Timestamp string `json:"timestamp"`
}

type Conversation struct {
	Id int64 `json:"id"`
	User int64 `json:"user"`
}

type Message struct {
	Id int64 `json:"id"`
	Text string `json:"text"`
	Image *os.File `json:"image"`
	Sender int64 `json:"sender"`
	Conversation int64 `json:"conversation"`
	Timestamp string `json:"timestamp"`
	IsForwarded bool `json:"isForwarded"`
	ForwardedMessage *int64 `json:"forwardedMessage"`
}

type Status struct {
	Message int64 `json:"message"`
	User int64 `json:"user"`
	Info string `json:"info"`
}

type Reaction struct {
	Emoji string `json:"emoji"`
	Message int64 `json:"message"`
	User int64 `json:"user"`
}
