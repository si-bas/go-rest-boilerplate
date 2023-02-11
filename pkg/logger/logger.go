package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/si-bas/go-rest-boilerplate/config"
	logCtx "github.com/si-bas/go-rest-boilerplate/pkg/logger/context"
	"github.com/si-bas/go-rest-boilerplate/pkg/logger/tag"
)

// Logger is logger singleton
var Logger StandardLogger

// StandardLogger is standard logger struct type
type StandardLogger struct {
	ZeroLogger zerolog.Logger
}

// InitLogger to initiate Logger
func InitLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Config.App.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	Logger.ZeroLogger = zerolog.New(zerolog.SyncWriter(os.Stdout)).With().
		Caller().
		Timestamp().
		Str("app_name", config.Config.App.Name).
		Logger()
}

func writeZeroLog(ev *zerolog.Event, tags ...tag.Tag) *zerolog.Event {
	for _, t := range tags {
		ev = ev.Str(t.Key, t.Value)
	}
	return ev
}

func getStackTrace(from, to int) (traces []string) {
	for from < to {
		if _, f, l, ok := runtime.Caller(from); ok {
			traces = append(traces, fmt.Sprintf("called from %s:%d", f, l))
		}
		from++
	}
	return traces
}

// Debug to provide global zerolog log Debug
func Debug(ctx context.Context, msg string, tags ...tag.Tag) {
	writeZeroLog(
		Logger.ZeroLogger.Debug().CallerSkipFrame(1),
		append(tags, logCtx.GetAllLoggingTagInTagStr(ctx)...)...).
		Msg(msg)
}

// Info to provide global zerolog log Info
func Info(ctx context.Context, msg string, tags ...tag.Tag) {
	writeZeroLog(
		Logger.ZeroLogger.Info().CallerSkipFrame(1),
		append(tags, logCtx.GetAllLoggingTagInTagStr(ctx)...)...).
		Msg(msg)
}

// Warn to provide global zerolog log Warn
func Warn(ctx context.Context, msg string, tags ...tag.Tag) {
	writeZeroLog(
		Logger.ZeroLogger.Warn().CallerSkipFrame(1),
		append(tags, logCtx.GetAllLoggingTagInTagStr(ctx)...)...,
	).Msg(msg)
}

// Error to provide global zerolog log Error
func Error(ctx context.Context, msg string, err error, tags ...tag.Tag) {
	ev := writeZeroLog(
		Logger.ZeroLogger.Error().CallerSkipFrame(1),
		append(tags, logCtx.GetAllLoggingTagInTagStr(ctx)...)...,
	)
	if err != nil {
		ev.Err(err)
	}
	ev.Strs("stack_trace", getStackTrace(2, 4)).Msg(msg)
}

// Fatal to provide global zerolog log Fatal
func Fatal(ctx context.Context, msg string, tags ...tag.Tag) {
	writeZeroLog(
		Logger.ZeroLogger.Fatal().CallerSkipFrame(1),
		append(tags, logCtx.GetAllLoggingTagInTagStr(ctx)...)...,
	).Msg(msg)
}
