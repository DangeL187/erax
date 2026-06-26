package main

import (
	"errors"
	"io"
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

func eraxWorkflow() string {
	err := erax.New(errMsgRoot)
	err = erax.WithMeta(err, "failed to create user",
		erax.F("code", errCode),
		erax.F("info", errInfo),
		erax.F("user_error", errUserError),
	)

	err = erax.WrapWithErrors(err, errMsgWrap, erax.New("random error"))

	return erax.FormatToJSONString(err)
}

func zapWorkflow(logger *zap.Logger) {
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

func BenchmarkEraxWorkflow(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res := eraxWorkflow()
		if len(res) == 0 {
			b.Fatal("empty result")
		}
	}
}

func BenchmarkZapWorkflow(b *testing.B) {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(io.Discard),
		zap.ErrorLevel,
	)
	logger := zap.New(core)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		zapWorkflow(logger)
	}
}
