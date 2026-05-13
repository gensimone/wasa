package database

import "fmt"

type EmojiCode string

const (
	Like  EmojiCode = "like"
	Love  EmojiCode = "love"
	Laugh EmojiCode = "laugh"
	Sad   EmojiCode = "sad"
	Angry EmojiCode = "angry"
	Thumb EmojiCode = "thumb"
)

// Returns true if the provided emoji code is a valid EmojiCode, otherwise false.
func IsValidEmojiCode(emojiCode EmojiCode) error {
	switch emojiCode {
	case Like, Love, Laugh, Sad, Angry, Thumb:
		return fmt.Errorf("Invalid emoji code: %s", emojiCode)
	default:
		return nil
	}
}
