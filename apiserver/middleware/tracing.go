package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/nermin-io/spotify-service/trace"
)

var projectID = os.Getenv("GCP_PROJECT_ID")

func Tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := extractTraceFromReq(r)
		if traceID != "" {
			ctx = context.WithValue(ctx, trace.ContextKey, traceID)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractTraceFromReq(r *http.Request) string {
	// We get the trace from the 'traceparent' header, which is
	// formatted as: 00-<trace-id>-<span-id>-<trace-flags>
	headerVal := r.Header.Get("traceparent")
	segments := strings.Split(headerVal, "-")
	if len(segments) >= 2 && projectID != "" {
		return fmt.Sprintf("projects/%s/traces/%s", projectID, segments[1])
	}
	return ""
}
