package model

import (
	"time"
	"unicode/utf8"
)

// Snippet — сохранённый текстовый фрагмент.
type Snippet struct {
	ID        string
	Content   string
	CreatedAt time.Time
	ExpiresAt *time.Time // nil — хранится без ограничения срока
}

// Length возвращает число символов в тексте записи.
func (s Snippet) Length() int {
	return utf8.RuneCountInString(s.Content)
}
