package log

import (
	"context"
	"io"
)

type (
	Log interface {
		Debug(s string)
		DebugWithContext(ctx context.Context, s string)
		Info(s string)
		InfoWithContext(ctx context.Context, s string)
		Warn(s string)
		WarnWithContext(ctx context.Context, s string)
		Err(s string)
		ErrWithContext(ctx context.Context, s string)
		Panic(s string)
		PanicWithContext(ctx context.Context, s string)
	}

	Logger struct {
		output io.Writer
		level  Level
	}

	Level uint32
)

const (
	Debug Level = iota
	Info
	Warn
	Err
	Panic
)

func NewLogger(output io.Writer, level Level)
