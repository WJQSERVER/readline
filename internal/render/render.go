package render

import (
	"fmt"
	"io"

	"github.com/mattn/go-runewidth"
	"github.com/WJQSERVER/readline/internal/buffer"
)

type Renderer struct {
	out          io.Writer
	prompt       string
	lastWidth    int
}

func NewRenderer(out io.Writer) *Renderer {
	return &Renderer{
		out: out,
	}
}

func (r *Renderer) SetPrompt(prompt string) {
	r.prompt = prompt
}

func (r *Renderer) Refresh(b *buffer.Buffer) error {
	currentWidth := b.FullWidth()
	cursorPos := b.DisplayWidth(b.Cursor())
	promptWidth := runewidth.StringWidth(r.prompt)

	// Basic redraw: carriage return, print prompt + content, clear to EOL
	fmt.Fprintf(r.out, "\r%s%s\x1b[K", r.prompt, b.String())

	// Move cursor to correct position
	// Use absolute horizontal position if supported, or move from start
	// \x1b[G is move to column (1-based)
	fmt.Fprintf(r.out, "\r\x1b[%dC", promptWidth+cursorPos)

	r.lastWidth = currentWidth
	return nil
}

func (r *Renderer) ClearLine() {
	fmt.Fprintf(r.out, "\r\x1b[K")
}

func (r *Renderer) NewLine() {
	fmt.Fprintf(r.out, "\r\n")
}
