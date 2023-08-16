package gorm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	traceStr      = "[%.3fms] [rows:%v] %s"
	traceWarnStr  = "SLOW SQL >= %v [%.3fms] [rows:%v] %s"
	slowThreshold = 200 * time.Millisecond
)

type gormLogger struct {
	log *slog.Logger
}

func (l gormLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return l
}

func (l gormLogger) Info(_ context.Context, msg string, data ...any) {
	l.log.Info(fmt.Sprintf(msg, data), "file", utils.FileWithLineNum())
}

// Warn print warn messages
func (l gormLogger) Warn(_ context.Context, msg string, data ...any) {
	l.log.Warn(fmt.Sprintf(msg, data), "file", utils.FileWithLineNum())
}

// Error print error messages
func (l gormLogger) Error(_ context.Context, msg string, data ...any) {
	l.log.Error(fmt.Sprintf(msg, data), "file", utils.FileWithLineNum())
}

// Trace print sql message
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	switch {
	case err != nil && l.log.Enabled(ctx, slog.LevelError) && !errors.Is(err, gorm.ErrRecordNotFound):
		sql, rows := fc()
		var rowsAny any
		if rows == -1 {
			rowsAny = "-"
		} else {
			rowsAny = rows
		}
		l.log.Error(fmt.Sprintf(traceStr, float64(elapsed.Nanoseconds())/1e6, rowsAny, sql), "error", err, "file", utils.FileWithLineNum())
	case elapsed > slowThreshold && l.log.Enabled(ctx, slog.LevelWarn):
		sql, rows := fc()
		var rowsAny any
		if rows == -1 {
			rowsAny = "-"
		} else {
			rowsAny = rows
		}
		l.log.Warn(fmt.Sprintf(traceWarnStr, slowThreshold, float64(elapsed.Nanoseconds())/1e6, rowsAny, sql), "file", utils.FileWithLineNum())
	case l.log.Enabled(ctx, slog.LevelInfo):
		sql, rows := fc()
		var rowsAny any
		if rows == -1 {
			rowsAny = "-"
		} else {
			rowsAny = rows
		}
		l.log.Info(fmt.Sprintf(traceStr, float64(elapsed.Nanoseconds())/1e6, rowsAny, sql), "file", utils.FileWithLineNum())
	}
}

func (gormLogger) ParamsFilter(_ context.Context, sql string, params ...any) (string, []any) {
	return sql, params
}
