package database

type Emoji string

const (
	Like  Emoji = "like"
	Love  Emoji = "love"
	Laugh Emoji = "laugh"
	Sad   Emoji = "sad"
	Angry Emoji = "angry"
	Thumb Emoji = "thumb"
)

var validEmoji = map[Emoji]struct{}{
	Like:  {},
	Love:  {},
	Laugh: {},
	Sad:   {},
	Angry: {},
	Thumb: {},
}

func IsValidEmoji(emoji Emoji) bool {
	_, ok := validEmoji[emoji]
	return ok
}
