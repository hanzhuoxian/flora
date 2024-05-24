package log

import "context"

type key int

const (
	logContextKey key = iota
)

func (z *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, z)
}

func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return WithName("Unknown-Context")
	}
	logger := ctx.Value(logContextKey)
	if logger != nil {
		return logger.(Logger)
	}

	return WithName("Unknown-Context")
}
