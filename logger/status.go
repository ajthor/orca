package logger

// import (
//   "fmt"
// )

// func (m *Logger) UpdateStatus(offset int, status string) {
//   logLength := len(m.log)
//
//   line := m.log[logLength-offset]
//   lineLength := len(line)
//
//   statusLength := len(status)
//
//   m.MoveToBeginning()
//   m.MoveUp(offset)
//
//   m.MoveForward(lineLength)
//
//   fmt.Printf(" %s\n", status)
//
//   m.MoveDown(offset-1)
// }
