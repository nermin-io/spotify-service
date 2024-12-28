package trace

import "context"

const (
	ContextKey = "trace-id"
)

func FromContext(ctx context.Context) string {
	traceID, ok := ctx.Value(ContextKey).(string)
	if !ok {
		return ""
	}
	return traceID
}
