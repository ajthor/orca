package logger

import (
  "fmt"
)
//
// func (sm *StatusLogger) Printf(msg string, args ...interface{}) {
//   spmsg := fmt.Sprintf(msg, args...)
//   sm.len = len(msg)
//   sm.m.Printf(spmsg)
// }

// Clear the line to the end.
func (m *Logger) Clear() {
  fmt.Printf("\033[K")
}

// Clear the whole line.
func (m *Logger) ClearLine() {
  m.MoveToBeginning()
  m.Clear()
}

// Move the cursor to a specific X/Y. Should not be used, except in certain
// circumstances.
func (m *Logger) Move(x int, y int) {
	fmt.Printf("\033[%d;%dH", x, y)
}

func (m *Logger) MoveUp(n int) {
  fmt.Printf("\033[%dA", n)
}

func (m *Logger) MoveDown(n int) {
  fmt.Printf("\033[%dB", n)
}

func (m *Logger) MoveForward(n int) {
  fmt.Printf("\033[%dC", n)
}

func (m *Logger) MoveBackward(n int) {
  fmt.Printf("\033[%dD", n)
}

func (m *Logger) MoveToBeginning() {
  fmt.Printf("\r")
}
