package buffer

import (
	"github.com/mattn/go-runewidth"
)

type Buffer struct {
	data   []rune
	cursor int // character position
}

func NewBuffer() *Buffer {
	return &Buffer{
		data: make([]rune, 0),
	}
}

func (b *Buffer) Insert(r rune) {
	b.data = append(b.data, 0)
	copy(b.data[b.cursor+1:], b.data[b.cursor:])
	b.data[b.cursor] = r
	b.cursor++
}

func (b *Buffer) Delete() {
	if b.cursor < len(b.data) {
		b.data = append(b.data[:b.cursor], b.data[b.cursor+1:]...)
	}
}

func (b *Buffer) Backspace() {
	if b.cursor > 0 {
		b.cursor--
		b.data = append(b.data[:b.cursor], b.data[b.cursor+1:]...)
	}
}

func (b *Buffer) MoveLeft() {
	if b.cursor > 0 {
		b.cursor--
	}
}

func (b *Buffer) MoveRight() {
	if b.cursor < len(b.data) {
		b.cursor++
	}
}

func (b *Buffer) MoveHome() {
	b.cursor = 0
}

func (b *Buffer) MoveEnd() {
	b.cursor = len(b.data)
}

func (b *Buffer) String() string {
	return string(b.data)
}

func (b *Buffer) Runes() []rune {
	return b.data
}

func (b *Buffer) Cursor() int {
	return b.cursor
}

func (b *Buffer) SetContent(s string) {
	b.data = []rune(s)
	b.cursor = len(b.data)
}

func (b *Buffer) Clear() {
	b.data = b.data[:0]
	b.cursor = 0
}

// DisplayWidth returns the visual width of the buffer up to a certain point
func (b *Buffer) DisplayWidth(limit int) int {
	if limit > len(b.data) {
		limit = len(b.data)
	}
	return runewidth.StringWidth(string(b.data[:limit]))
}

func (b *Buffer) FullWidth() int {
	return runewidth.StringWidth(string(b.data))
}
