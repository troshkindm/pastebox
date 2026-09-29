package model

import "time"

// Snippet — сохранённый текстовый фрагмент.
type Snippet struct {
	ID        string
	Content   string
	CreatedAt time.Time
}
