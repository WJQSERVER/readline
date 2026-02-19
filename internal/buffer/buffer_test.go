package buffer

import "testing"

func TestBuffer(t *testing.T) {
	b := NewBuffer()
	b.Insert('h')
	b.Insert('e')
	b.Insert('l')
	b.Insert('l')
	b.Insert('o')
	if b.String() != "hello" {
		t.Errorf("expected hello, got %s", b.String())
	}
	if b.Cursor() != 5 {
		t.Errorf("expected cursor 5, got %d", b.Cursor())
	}

	b.MoveLeft()
	b.Insert(' ')
	if b.String() != "hell o" {
		t.Errorf("expected 'hell o', got '%s'", b.String())
	}

	b.Backspace()
	if b.String() != "hello" {
		t.Errorf("expected hello, got %s", b.String())
	}
}

func TestUnicode(t *testing.T) {
	b := NewBuffer()
	b.Insert('你')
	b.Insert('好')
	if b.String() != "你好" {
		t.Errorf("expected 你好, got %s", b.String())
	}
	if b.FullWidth() != 4 {
		t.Errorf("expected width 4, got %d", b.FullWidth())
	}
}
