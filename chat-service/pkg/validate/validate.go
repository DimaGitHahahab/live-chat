package validate

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

const (
	MinNicknameLength = 3
	MaxNicknameLength = 20
	MaxMessageLength  = 120
)

var (
	ErrShortNickname  = errors.New(fmt.Sprintf("nickname is too short, must be >= %d symbols", MinNicknameLength))
	ErrLongNickname   = errors.New(fmt.Sprintf("nickname is too long, must be <= %d symbols", MaxNicknameLength))
	ErrEmptyMessage   = errors.New("message cannot be empty")
	ErrTooLongMessage = errors.New(fmt.Sprintf("message is too long, must be <= %d symbols", MaxMessageLength))
)

func Name(n string) error {
	length := utf8.RuneCountInString(n)
	if length < MinNicknameLength {
		return ErrShortNickname
	}
	if length > MaxNicknameLength {
		return ErrLongNickname
	}
	return nil
}

func Message(m string) error {
	length := utf8.RuneCountInString(m)
	if length == 0 {
		return ErrEmptyMessage
	}
	if length > MaxMessageLength {
		return ErrTooLongMessage
	}

	return nil
}
