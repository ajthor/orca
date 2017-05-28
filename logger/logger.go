package logger

import (
  "fmt"
  "log"

  . "github.com/logrusorgru/aurora"
)

type Logger struct {
  opts LoggerOptions
  log []Message
}

type LoggerOptions struct {
  Verbose int
  Colors bool
}

type Message interface {
  fmt.Stringer
}

type message struct {
  spmsg interface{}
}

func (m message) String() string {
  return fmt.Sprint(m.spmsg)
}

func (m LogMessage) Message() string { return m.spmsg }

func NewLogger(opts LoggerOptions) *Logger {
  // For now, everything is verbose.
  opts.Verbose = 3
  opts.Colors = true

  return &Logger{opts: opts}
}

func (m *Logger) Printf(msg string, args ...interface{}) {
  spmsg := fmt.Sprintf(msg, args...)
  // m.log = append(m.log, spmsg)
  fmt.Printf(spmsg)
}

func (m *Logger) Info(msg interface{}, args ...interface{}) LogMessage {
  if val, ok := msg.(value); ok {
    val.spmsg = fmt.Sprintf("%s%s", Green("[INFO]  "), val)
    return val
  }

  return {spmsg}
  // alert := Green("[INFO]  ")
  // if m.opts.Verbose >= 3 {
  //   m.Printf("%s%s", alert, msg, args...)
  // }
  spmsg := fmt.Sprintf("%s%s", Green("[INFO]  "), msg)
  if m.opts.Verbose >= 3 {
    m.Printf(spmsg, args...)
  }
}

func (m *Logger) Warn(msg string, args ...interface{}) {
  // alert := Brown("[WARN]  ")
  // if m.opts.Verbose >= 2 {
  //   m.Printf("%s%s", alert, msg, args...)
  // }
  spmsg := fmt.Sprintf("%s%s", Brown("[WARN]  "), msg)
  if m.opts.Verbose >= 2 {
    m.Printf(spmsg, args...)
  }
}

func (m *Logger) Error(msg string, args ...interface{}) {
  spmsg := fmt.Sprintf("%s%s", Red("[ERROR] "), msg)
  if m.opts.Verbose >= 1 {
    m.Printf(spmsg, args...)
  }
}

func (m *Logger) Fatal(msg interface{})  {
  alert := Red("[ERROR] ")
  switch t := msg.(type) {
  case string:
    spmsg := fmt.Sprintf("%s%s", alert, t)
    m.Printf(spmsg)
    log.Fatal(t)
  case error:
    log.Fatal(t)
  }
}
