package readline

import (
	"errors"
	"io"
	"github.com/WJQSERVER/readline/internal/buffer"
	"github.com/WJQSERVER/readline/internal/input"
	"github.com/WJQSERVER/readline/internal/render"
	"github.com/WJQSERVER/readline/internal/term"
)

var (
	ErrInterrupt = errors.New("interrupt")
	ErrEOF       = io.EOF
)

type Instance struct {
	cfg      *Config
	terminal term.Terminal
	buffer   *buffer.Buffer
	renderer *render.Renderer
	parser   *input.Parser

	historyIdx int // -1 means current line
	tempBuffer string
}

func NewInstance(cfg *Config) (*Instance, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	cfg.Init()

	t, err := term.NewTerminal(cfg.Stdin, cfg.Stdout)
	if err != nil {
		return nil, err
	}

	return &Instance{
		cfg:      cfg,
		terminal: t,
		buffer:   buffer.NewBuffer(),
		renderer: render.NewRenderer(t),
		parser:   input.NewParser(t),
		historyIdx: -1,
	}, nil
}

func (i *Instance) Readline() (string, error) {
	restore, err := i.terminal.SetRaw()
	if err != nil {
		return "", err
	}
	defer restore()

	i.buffer.Clear()
	i.historyIdx = -1
	i.renderer.SetPrompt(i.cfg.Prompt)
	i.renderer.Refresh(i.buffer)

	for {
		ev, err := i.parser.NextEvent()
		if err != nil {
			return "", err
		}

		switch ev.Key {
		case input.KeyEnter:
			line := i.buffer.String()
			i.renderer.NewLine()
			i.cfg.History.Append(line)
			return line, nil
		case input.KeyBackspace:
			i.buffer.Backspace()
		case input.KeyDelete:
			i.buffer.Delete()
		case input.KeyLeft:
			i.buffer.MoveLeft()
		case input.KeyRight:
			i.buffer.MoveRight()
		case input.KeyCtrlLeft:
			i.buffer.MoveWordLeft()
		case input.KeyCtrlRight:
			i.buffer.MoveWordRight()
		case input.KeyHome:
			i.buffer.MoveHome()
		case input.KeyEnd:
			i.buffer.MoveEnd()
		case input.KeyUp:
			i.handleHistory(true)
		case input.KeyDown:
			i.handleHistory(false)
		case input.KeyCtrlC:
			i.renderer.NewLine()
			return "", ErrInterrupt
		case input.KeyCtrlD:
			if i.buffer.String() == "" {
				return "", ErrEOF
			}
			i.buffer.Delete()
		case input.KeyCtrlL:
			i.renderer.Refresh(i.buffer) // Just refresh, clear logic is in renderer
		case input.KeyRune:
			i.buffer.Insert(ev.Rune)
		case input.KeyTab:
			i.handleCompletion()
		}

		i.renderer.Refresh(i.buffer)
	}
}

func (i *Instance) handleHistory(up bool) {
	if up {
		if i.historyIdx == -1 {
			i.tempBuffer = i.buffer.String()
			i.historyIdx = i.cfg.History.Len() - 1
		} else if i.historyIdx > 0 {
			i.historyIdx--
		} else {
			return
		}
	} else {
		if i.historyIdx == -1 {
			return
		}
		i.historyIdx++
		if i.historyIdx >= i.cfg.History.Len() {
			i.historyIdx = -1
			i.buffer.SetContent(i.tempBuffer)
			return
		}
	}

	if line, ok := i.cfg.History.Get(i.historyIdx); ok {
		i.buffer.SetContent(line)
	}
}

func (i *Instance) handleCompletion() {
	if i.cfg.Completer == nil {
		return
	}
	candidates, length := i.cfg.Completer.Do(i.buffer.Runes(), i.buffer.Cursor())
	if len(candidates) == 1 {
		// Single match, just insert it
		suffix := candidates[0][length:]
		for _, r := range suffix {
			i.buffer.Insert(r)
		}
	}
}

func (i *Instance) Close() error {
	return nil
}
