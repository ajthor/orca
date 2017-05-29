package logger

import (
  "fmt"
)

// The functions below provide simple terminal manipulation to move the cursor.
func (l *Logger) Clear() {
  fmt.Printf("\033[K")
}

func (l *Logger) ClearScreen() {
  fmt.Printf("\033[2J")
}

// Should not be used, except in certain circumstances.
func (l *Logger) Move(x int, y int) {
	fmt.Printf("\033[%d;%dH", x, y)
}

func (l *Logger) MoveUp(n int) {
  fmt.Printf("\033[%dA", n)
}

func (l *Logger) MoveDown(n int) {
  fmt.Printf("\033[%dB", n)
}

func (l *Logger) MoveForward(n int) {
  fmt.Printf("\033[%dC", n)
}

func (l *Logger) MoveBackward(n int) {
  fmt.Printf("\033[%dD", n)
}

func (l *Logger) MoveToBeginning() {
  fmt.Printf("\r")
}

func (l *Logger) SaveCursorPosition() {
  fmt.Printf("\033[s")
}

func (l *Logger) RestoreCursorPosition() {
  fmt.Printf("\033[u")
}

// The functions below are for modifying the history of the logger.
func (l *Logger) modify(index int, msg *Message) {
  histLen := len(l.history)
  var offset int

  // offset = histLen - offset
  if index < 0 {
    offset = -1*index
  } else {
    offset = histLen - index
  }

  // fmt.Printf("Len: %d, Offset: %d, Index: %d", histLen, offset, index)
  // if index > histLen {
  //   return errors.New("Offset greater than history length.")
  // }

  l.SaveCursorPosition()
  l.MoveToBeginning()
  l.MoveUp(offset)

  l.Clear()

  rmsg := l.format(msg)
  fmt.Printf(rmsg)

  l.RestoreCursorPosition()
}

func (l *Logger) Modify(index int, lvl Level, args ...interface{}) (msg *Message) {
  msg = NewMessage(lvl, nil, args...)
  l.modify(index, msg)
  return msg
}

func (l *Logger) Modifyf(index int, lvl Level, format interface{}, args ...interface{}) (msg *Message) {
  msg = NewMessage(lvl, format, args...)
  l.modify(index, msg)
  return msg
}

// The functions below are accessible by the 'std' logger, defined in
// `logger.go`.
func Modify(index int, lvl Level, args ...interface{}) (msg *Message) {
  return std.Modify(index, lvl, args...)
}

func Modifyf(index int, lvl Level, format interface{}, args ...interface{}) (msg *Message) {
  return std.Modifyf(index, lvl, format, args...)
}
