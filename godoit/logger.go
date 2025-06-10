package godoit

import "context"

// LogIt is an interface to determine what messages should be logged at what level, this does not actually indicate what
// level a log should be emitted at but rather a suggestion of how the log should be logged
// TODO turn into a struct that the functions can be changed for
type LogIt interface {
	// TraceLog is used to indicate this msg should be logged at the trace level
	TraceLog(ctx context.Context, msg string)

	// DebugLog is used to indicate this msg should be logged at the debug level
	DebugLog(ctx context.Context, msg string)

	// InfoLog is used to to indicate this msg should be logged at the info level
	InfoLog(ctx context.Context, msg string)

	// WarnLog is used to indicate this msg should be logged at the warn level
	WarnLog(ctx context.Context, msg string, err error)

	// ErrorLog is used to indicate this msg should be logged at the error level
	ErrorLog(ctx context.Context, msg string, err error)

	// FatalLog is used to indicate this msg should be logged at the fatal level
	FatalLog(ctx context.Context, msg string, err error)
}

type DefaultLogger struct {
	ctx context.Context
}
