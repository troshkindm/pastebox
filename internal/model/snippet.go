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
}

// Length возвращает число символов в тексте записи.
func (s Snippet) Length() int {
	return utf8.RuneCountInString(s.Content)
}
