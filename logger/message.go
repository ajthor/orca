package logger

import (
  "fmt"
  "strings"
)

type message interface {
  fmt.Stringer
}

type Message struct {
  lvl Level
  prefix string
  format *string
  args []interface{}
}

func (m *Message) String() string {
  // TODO: Add prefix here. Likely, we will need an intermediate variable to
  // append the prefix, and then we will pass it into the Sprintf functions
  // below.

  if m.format != nil {
    return fmt.Sprintf(*m.format, m.args...)
  } else {
    return fmt.Sprint(m.args...)
  }
}

func (m *Message) append(spmsg string) {
  if m.format != nil {
    spmsg = strings.TrimRight(*m.format, "\n") + spmsg
    m.format = &spmsg
  } else {
    m.args = append(m.args, spmsg)
  }
}

func NewMessage(lvl Level, format interface{}, args ...interface{}) *Message {
  msg := &Message{
    lvl: lvl,
    args: args,
  }

  switch f := format.(type) {
  case string:
    msg.format = &f
  case *string:
    msg.format = f
  }

  return msg
}
