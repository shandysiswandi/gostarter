package requestid

import "context"

type requestIDKey struct{}

func Set(ctx context.Context, rID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, rID)
}

func Get(ctx context.Context) string {
	if val := ctx.Value(requestIDKey{}); val != nil {
		if requestID, ok := val.(string); ok {
			return requestID
		}
	}

	return ""
}
