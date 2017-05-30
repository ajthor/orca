package logger

import (
  "fmt"
  "os"
  "strings"
  "sync"

  . "github.com/logrusorgru/aurora"
)

type Level int

const (
  FATAL Level = iota
  ERROR
  WARNING
  INFO
  DEBUG
)

type logger interface {
  Fatal()
  Fatalf()
  Error()
  Errorf()
  Warn()
  Warnf()
  Info()
  Infof()
  Debug()
  Debugf()
}

type Logger struct {
  mu sync.Mutex
  logLevel Level
  history []*Message
}

// Creates a new instance of the logger.
func New(lvl Level) *Logger {
  // For now, everything is verbose.
  return &Logger{logLevel: lvl}
}

func (l *Logger) SetLogLevel(lvl Level) {
  l.mu.Lock()
  defer l.mu.Unlock()

  l.logLevel = lvl
}

func (l *Logger) format(msg *Message) string {
  spmsg := msg.String()

  // Trim any trailing newlines. We want the log to be clean.
  spmsg = strings.TrimRight(spmsg, "\n")
  // Append a single newline to the end of the formatted string.
  if last := len(spmsg) - 1; spmsg[last] != 10 && spmsg[last] != 109 {
    spmsg = spmsg + "\n"
  }

  switch msg.lvl {
  case FATAL:
    return fmt.Sprintf("%s%s", Red("[FATAL]  "), spmsg)
  case ERROR:
    return fmt.Sprintf("%s%s", Red("[ERROR]  "), spmsg)
  case WARNING:
    return fmt.Sprintf("%s%s", Brown("[WARN]   "), spmsg)
  case INFO:
    return fmt.Sprintf("%s%s", Green("[INFO]   "), spmsg)
  case DEBUG:
    return fmt.Sprintf("%s%s", Cyan("[DEBUG]  "), spmsg)
  // By default, any messages that do not fall into one of the other log levels
  // are output as 'DEBUG'.
  default:
    return fmt.Sprintf("%s%s", Cyan("[DEBUG]  "), spmsg)
  }
}

func (l *Logger) log(lvl Level, format *string, args ...interface{}) *Message {
  msg := NewMessage(lvl, format, args...)

  l.mu.Lock()
  defer l.mu.Unlock()
  l.history = append(l.history, msg)

  // Could implement different types of logging here. Output to a file, for
  // example.
  if lvl <= l.logLevel {
    rmsg := l.format(msg)
    fmt.Print(rmsg)
  }

  return msg
}

func (l *Logger) Fatal(args ...interface{}) (msg *Message) {
  msg = l.log(FATAL, nil, args...)
  os.Exit(1)
  return
}

func (l *Logger) Fatalf(format string, args ...interface{}) (msg *Message) {
  msg = l.log(FATAL, &format, args...)
  os.Exit(1)
  return
}

func (l *Logger) Error(args ...interface{}) (msg *Message) {
  msg = l.log(ERROR, nil, args...)
  panic(fmt.Sprint(args...))
  return
}

func (l *Logger) Errorf(format string, args ...interface{}) (msg *Message) {
  msg = l.log(ERROR, &format, args...)
  panic(fmt.Sprintf(format, args...))
  return
}

func (l *Logger) Warn(args ...interface{}) (msg *Message) {
  msg = l.log(WARNING, nil, args...)
  return
}

func (l *Logger) Warnf(format string, args ...interface{}) (msg *Message) {
  msg = l.log(WARNING, &format, args...)
  return
}

func (l *Logger) Info(args ...interface{}) (msg *Message) {
  msg = l.log(INFO, nil, args...)
  return
}

func (l *Logger) Infof(format string, args ...interface{}) (msg *Message) {
  msg = l.log(INFO, &format, args...)
  return
}

func (l *Logger) Debug(args ...interface{}) (msg *Message) {
  msg = l.log(DEBUG, nil, args...)
  return
}

func (l *Logger) Debugf(format string, args ...interface{}) (msg *Message) {
  msg = l.log(DEBUG, &format, args...)
  return
}

// Global declaration of a logger. This logger is available through the exposed
// methods founds below.
var std = New(DEBUG)

// Methods for calling the default 'std' logger. These are exported and
// available for use by importing the package in another file. This effectively
// makes the logger global.
func Fatal(args ...interface{}) (msg *Message) {
  return std.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) (msg *Message) {
  return std.Fatalf(format, args...)
}

func Error(args ...interface{}) (msg *Message) {
  return std.Error(args...)
}

func Errorf(format string, args ...interface{}) (msg *Message) {
  return std.Errorf(format, args...)
}

func Warn(args ...interface{}) (msg *Message) {
  return std.Warn(args...)
}

func Warnf(format string, args ...interface{}) (msg *Message) {
  return std.Warnf(format, args...)
}

func Info(args ...interface{}) (msg *Message) {
  return std.Info(args...)
}

func Infof(format string, args ...interface{}) (msg *Message) {
  return std.Infof(format, args...)
}

func Debug(args ...interface{}) (msg *Message) {
  return std.Debug(args...)
}

func Debugf(format string, args ...interface{}) (msg *Message) {
  return std.Debugf(format, args...)
}
