package slog

import (
	"context"
	"github.com/m2c/kiplestar/commons"
	"strconv"
	"sync/atomic"
)


type Span struct {
	ParentSpanID int32
	SubSpanID int32
	NextSpanID int32
}

func GetSubSpanID(ctx context.Context, s *Span) string {
	tempVal := &s.SubSpanID
	s.SubSpanID = atomic.AddInt32(tempVal, 1)
	context.WithValue(ctx, commons.X_SPAN_ID, s)
	return strconv.Itoa(int(s.ParentSpanID)) + "." + strconv.Itoa(int(s.SubSpanID))
}

func GetNextSpanID(ctx context.Context, s *Span) int32 {
	tempVal := &s.NextSpanID
	s.NextSpanID = atomic.AddInt32(tempVal, 1)
	context.WithValue(ctx, commons.X_SPAN_ID, s)
	return s.NextSpanID
}
