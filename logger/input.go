package logger

import (
  "bufio"
  "fmt"
  "os"
  "strings"

  . "github.com/logrusorgru/aurora"
)

type Response int

const (
  NONE Response = iota
  YES
  NO
  CANCEL
)

const (
  DEFAULT Level = iota
  YESNO
  YESNOCANCEL
)

func (l *Logger) formatPrompt(msg *Message) string {
  spmsg := msg.String()

  // Trim any trailing newlines. We want the log to be clean.
  spmsg = strings.TrimRight(spmsg, "\n")
  spmsg = strings.TrimRight(spmsg, ":")
  // Append a single newline to the end of the formatted string.
  // if last := len(spmsg) - 1; spmsg[last] != 10 && spmsg[last] != 109 {
  //   spmsg = spmsg + "\n"
  // }

  switch msg.lvl {
  case DEFAULT:
    return fmt.Sprintf("%s: ", spmsg)
  case YESNO:
    return fmt.Sprintf("%s [%s]: ", spmsg, Brown("Yn"))
  case YESNOCANCEL:
    return fmt.Sprintf("%s [%s]: ", spmsg, Brown("Ync"))
  default:
    return fmt.Sprintf("%s: ", spmsg)
  }
}

func (l *Logger) prompt(lvl Level, format *string, args ...interface{}) string {
  msg := NewMessage(lvl, format, args...)

  l.mu.Lock()
  defer l.mu.Unlock()
  l.history = append(l.history, msg)

  rmsg := l.formatPrompt(msg)
  fmt.Print(rmsg)

  res := ReadInput()

  return res
}

func (l *Logger) ShowInput(d string) {
  msg := l.history[len(l.history) - 1]

  // msg.args = append(msg.args, )
  // msg.args[last] =

  l.SaveCursorPosition()
  l.MoveToBeginning()
  l.MoveUp(1)

  l.Clear()

  rmsg := l.formatPrompt(msg)
  rmsg = strings.TrimRight(rmsg, "\n")
  fmt.Print(rmsg)
  fmt.Printf("%s\n", Cyan(d))

  l.RestoreCursorPosition()
}

func (l *Logger) FormatResponse(res string) Response {
  switch res {
  case "y", "Y", "yes", "Yes":
    return YES
  case "n", "N", "no", "No":
    return NO
  case "c", "C", "cancel", "Cancel":
    return CANCEL
  default:
    return NONE
  }
}

func (l *Logger) Prompt(lvl Level, args ...interface{}) string {
  return l.prompt(lvl, nil, args...)
}

func (l *Logger) Promptf(lvl Level, format string, args ...interface{}) string {
  return l.prompt(lvl, &format, args...)
}

func (l *Logger) ReadInput() string {
  scanner := bufio.NewScanner(os.Stdin)
  scanner.Scan()
  return scanner.Text()
}

// func (l *Logger) ReadInputMultiline() (lines []string) {
//   scanner := bufio.NewScanner(os.Stdin)
//   for scanner.Scan() {
//     lines = append(lines, scanner.Text())
//   }
//   return
// }

func Prompt(lvl Level, args ...interface{}) string {
  return std.Prompt(lvl, args...)
}

func Promptf(lvl Level, format string, args ...interface{}) string {
  return std.Promptf(lvl, format, args...)
}

func ShowInput(d string) {
  std.ShowInput(d)
}

func FormatResponse(res string) Response {
  return std.FormatResponse(res)
}

func ReadInput() string {
  return std.ReadInput()
}
