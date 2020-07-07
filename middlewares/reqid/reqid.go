package reqid

import (
	"context"
	"fmt"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Key to use when setting the request ID.
type ctxKeyRequestID string

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey ctxKeyRequestID = ""

//RequestID Creates a new request id
func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			u2 := uuid.NewV4()
			requestID = fmt.Sprintf("%s", u2)
		}
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

//GetReqID Creates a new request id
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
