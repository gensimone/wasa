package database

type EmojiCode string

const (
	Like  EmojiCode = "like"
	Love  EmojiCode = "love"
	Laugh EmojiCode = "laugh"
	Sad   EmojiCode = "sad"
	Angry EmojiCode = "angry"
	Thumb EmojiCode = "thumb"
)

type InvalidEmojiCodeError struct {
	Code EmojiCode
}

func (e *InvalidEmojiCodeError) Error() string {
	return "Invalid emoji code: " + string(e.Code)
}

func ValidateEmoji(emojiCode EmojiCode) error {
	switch emojiCode {
	case Like, Love, Laugh, Sad, Angry, Thumb:
		return &InvalidEmojiCodeError{Code: emojiCode}
	default:
		return nil
	}
}
