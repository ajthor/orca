package logger

// "fmt"

import (
  "strings"

  . "github.com/logrusorgru/aurora"
)

// type status interface {
//   fmt.Stringer
//
//   done()
// }
//
// type Status struct {
//   s string
//   msg *Message
// }

// The Done function is used to update the status of the previous statement by
// appending either "err" or "ok" to the end and changing the 'logLevel' label
// if the result is 'falsey'. Changing the label does not call the panic or
// os.Exit functions, but does provide a visual indicator of whether or not the
// previous command exited properly. Any panic or exit functions will need to
// be called manually.
// 'Falsey' in the context of the Done function is any value which would
// normally casue the program to stop. Falsey values in this context could be
// an error that is not nil or a false boolean value.
func (l *Logger) Done(args ...interface{}) (ok bool) {
  // Get the last message from history.
  msg := l.history[len(l.history) - 1]

  // Remove any trailing newlines from the string.
  last := len(msg.args) - 1
  msg.args[last] = strings.TrimRight(msg.args[last].(string), "\n")

  ok = true

  // Run through all of the arguments and determine if any of them are falsey.
  // If they are falsey, we change the `ok` variable to false.
  for _, arg := range args {
    switch t := arg.(type) {
    case error:
      ok = false
    case bool:
      if !t { ok = false }
    default:
      // fmt.Printf("type %T\n", t)
      // fmt.Println(arg)
    }
  }

  // Update the log message and modify the history.
  if ok {
    msg.lvl = INFO
    msg.args = append(msg.args, Green(" ok\n"))
  } else {
    msg.lvl = ERROR
    msg.args = append(msg.args, Red(" err\n"))
  }

  l.modify(-1, msg)

  return
}

func Done(args ...interface{}) bool {
  return std.Done(args...)
}
