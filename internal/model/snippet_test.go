package model_test

import (
	"testing"

	"github.com/troshkindm/pastebox/internal/model"
)

func TestSnippetLength(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"empty", "", 0},
		{"ascii", "hello", 5},
		{"cyrillic counts runes", "привет", 6},
		{"newline is a character", "a\nb", 3},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := model.Snippet{Content: tc.content}.Length()
			if got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}
