package logger

import (
  "github.com/sirupsen/logrus"
  "github.com/go-chi/chi/middleware"
  "net/http"
  "fmt"
  "time"
	reqid "golang-rest-api/middlewares/reqid"
  )

//StructuredLogger Logger type declaration
type StructuredLogger struct {
  Logger *logrus.Logger
}
//StructuredLoggerEntry Logger type structure
type StructuredLoggerEntry struct {
  Logger logrus.FieldLogger
}

//Logger FieldLogger type structure
func Logger(logger *logrus.Logger) func(next http.Handler) http.Handler {
  return middleware.RequestLogger(&StructuredLogger{logger})
}

//NewLogEntry Return entry for a new log line
func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
  entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
  scheme := "http"
  if r.TLS != nil {
    scheme = "https"
  }    
  logFields := logrus.Fields{
    "ts" : time.Now().UTC().Format(time.RFC1123),
    "req_id" : reqid.GetReqID(r.Context()),
    "http_scheme" : scheme,
    "http_proto" : r.Proto,
    "http_method" : r.Method,
    "remote_addr" : r.RemoteAddr,
    "user_agent" : r.UserAgent(),
    "uri" : fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI),
  }
  entry.Logger = entry.Logger.WithFields(logFields)
  return entry
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
l.Logger = l.Logger.WithFields(logrus.Fields{
  "resp_status": status, "resp_bytes_length": bytes,
  "resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
})
l.Logger.Info("Response served successfully")
}

//Panic func for defining log structure under appliation panic
func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
l.Logger = l.Logger.WithFields(logrus.Fields{
  "stack": string(stack),
  "panic": fmt.Sprintf("%+v", v),
})
}

//GetLogEntry func for reading a Log
func GetLogEntry(r *http.Request) logrus.FieldLogger {
  entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
  return entry.Logger
}

//LogEntrySetField Function to set field to an entry
func LogEntrySetField(r *http.Request, key string, value interface{}) {
  if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
    entry.Logger = entry.Logger.WithField(key, value)
  }
}

//LogEntrySetFields Function to set fields to an entry
func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
  if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
    entry.Logger = entry.Logger.WithFields(fields)
  }
}

//GetLoggerWithRID Return a logger object with req-id
func GetLoggerWithRID(logger *logrus.Logger, r  *http.Request) (requestLogger logrus.FieldLogger){
  requestLogger = logger.WithFields(logrus.Fields{
    "req-id": reqid.GetReqID(r.Context()),
  })
  return
}