package main

import (
	"errors"
	"io"
	"runtime"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/DangeL187/erax"
)

const (
	errCode      = "503"
	errInfo      = "This is a really\nreally long information."
	errUserError = "An account with this email already exists."
	errMsgRoot   = "email is already in use"
	errMsgWrap   = "failed to register\nbecause of ducks!"
)

var (
	globalResult string
	zapLogger    *zap.Logger
)

func init() {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(io.Discard),
		zap.ErrorLevel,
	)
	zapLogger = zap.New(core)
}

// ==========================================
// SCENARIO 1: LIGHT
// ==========================================

func eraxLight() string {
	err := erax.New(errMsgRoot)
	return erax.FormatToJSONString(err)
}

func zapLight(logger *zap.Logger) {
	logger.Error(errMsgRoot)
}

func Benchmark_Light_Erax(b *testing.B) {
	runtime.GC()
	b.ResetTimer()
	b.ReportAllocs()
	var res string
	for i := 0; i < b.N; i++ {
		res = eraxLight()
	}
	globalResult = res
}

func Benchmark_Light_Zap(b *testing.B) {
	runtime.GC()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		zapLight(zapLogger)
	}
}

// ==========================================
// SCENARIO 2: MEDIUM
// ==========================================

func eraxMedium() string {
	err := erax.New(errMsgRoot)
	err = erax.WithMeta(err, "failed to create user",
		erax.F("code", errCode),
		erax.F("info", errInfo),
		erax.F("user_error", errUserError),
	)
	err = erax.WrapWithErrors(err, errMsgWrap, errors.New("random error"))
	return erax.FormatToJSONString(err)
}

func zapMedium(logger *zap.Logger) {
	logger.Error(errMsgWrap,
		zap.Error(errors.New("random error")),
		zap.Dict("cause",
			zap.String("message", errMsgRoot),
			zap.Dict("meta",
				zap.String("code", errCode),
				zap.String("info", errInfo),
				zap.String("user_error", errUserError),
			),
		),
	)
}

func Benchmark_Medium_Erax(b *testing.B) {
	runtime.GC()
	b.ResetTimer()
	b.ReportAllocs()
	var res string
	for i := 0; i < b.N; i++ {
		res = eraxMedium()
	}
	globalResult = res
}

func Benchmark_Medium_Zap(b *testing.B) {
	runtime.GC()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		zapMedium(zapLogger)
	}
}

// ==========================================
// NEW SCENARIO 3: HEAVY
// ==========================================

func eraxHeavy() string {
	errDB := erax.New("db connection timeout")
	errDB = erax.WithMeta(errDB, "db_error", erax.F("host", "db-replica-1"), erax.F("retry", "3"))

	errAPI := erax.New("bad gateway")
	errAPI = erax.WithMeta(errAPI, "api_error", erax.F("provider", "stripe"), erax.F("status", "502"))

	errOrchestrator := erax.WrapWithErrors(errDB, "checkout failed", errAPI)
	errOrchestrator = erax.WithMeta(errOrchestrator, "orchestrator_error",
		erax.F("trace_id", "tx-9999"),
		erax.F("user_id", "usr-123"),
	)

	return erax.FormatToJSONString(errOrchestrator)
}

func zapHeavy(logger *zap.Logger) {
	dbLogger := logger.With(
		zap.String("context", "database"),
		zap.String("host", "db-replica-1"),
		zap.String("retry", "3"),
	)
	dbLogger.Debug("low-level db connection trace")

	apiLogger := logger.With(
		zap.String("context", "api_stripe"),
		zap.String("provider", "stripe"),
		zap.String("status", "502"),
	)
	apiLogger.Debug("low-level external api trace")

	orchestratorLogger := logger.With(
		zap.String("trace_id", "tx-9999"),
		zap.String("user_id", "usr-123"),
	)

	orchestratorLogger.Error(errMsgWrap,
		zap.String("db_cause", "db connection timeout"),
		zap.String("api_cause", "bad gateway"),
	)
}

func Benchmark_Heavy_Erax(b *testing.B) {
	runtime.GC()
	b.ResetTimer()
	b.ReportAllocs()
	var res string
	for i := 0; i < b.N; i++ {
		res = eraxHeavy()
	}
	globalResult = res
}

func Benchmark_Heavy_Zap(b *testing.B) {
	runtime.GC()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		zapHeavy(zapLogger)
	}
}
