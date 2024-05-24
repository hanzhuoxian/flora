package log

import (
	"context"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InfoLogger represents the ability to log non-error messages, at a particular verbosity.
type InfoLogger interface {
	Info(msg string, keysAndValues ...interface{})
	Infof(format string, args ...interface{})
	Enabled() bool
}

type Logger interface {
	InfoLogger

	Error(err error, msg string, keysAndValues ...interface{})
	Errorf(format string, args ...interface{})

	V(level int) InfoLogger

	Write(p []byte) (n int, err error)

	WithValues(keysAndValues ...interface{}) Logger

	WithName(name string) Logger

	WithContext(ctx context.Context) context.Context

	Flush()
}

var (
	_ Logger     = &zapLogger{}
	_ InfoLogger = &infoLogger{}
)

// noopInfoLogger is a no-op InfoLogger.
type noopInfoLogger struct{}

func (n noopInfoLogger) Enabled() bool                                 { return false }
func (n noopInfoLogger) Info(msg string, keysAndValues ...interface{}) {}
func (n noopInfoLogger) Infof(format string, args ...interface{})      {}

var disabledInfoLogger = &noopInfoLogger{}

type infoLogger struct {
	level zapcore.Level
	log   *zap.Logger
}

func (i *infoLogger) Enabled() bool {
	return true
}

func (i *infoLogger) Info(msg string, keysAndValues ...interface{}) {
	if checkEntry := i.log.Check(i.level, msg); checkEntry != nil {
		checkEntry.Write(handleFields(i.log, keysAndValues)...)
	}
}

func (i *infoLogger) Infof(format string, args ...interface{}) {
	i.log.Sugar().Infof(format, args...)
}

func Init(opts *Options) {
	options = opts
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if opts.EnableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zap.InfoLevel
	}

	loggerConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       false,
		DisableCaller:     !opts.EnableCaller,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         opts.Format,
		EncoderConfig:    encoderConfig,
		OutputPaths:      opts.OutputPaths,
		ErrorOutputPaths: opts.ErrorOutputPaths,
	}

	var err error
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	logger = &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
}

func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	if len(args) == 0 {
		return additional
	}
	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passwd to logr", zap.Any("field", args[i]))

			break
		}

		if i == len(args)-1 {
			l.DPanic("odd number of arguments passwd as key-value pairs for logging", zap.Any("ignored key", args[i]))

			break
		}

		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			l.DPanic(
				"non-string key argument passed to logging, ignoring all later arguments",
				zap.Any("invalid key", key),
			)

			break
		}
		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}

	return append(fields, additional...)
}

type zapLogger struct {
	zapLogger *zap.Logger
	infoLogger
}

var (
	logger  *zapLogger
	options *Options
)

func StdErrLogger() *log.Logger {
	if logger == nil {
		return nil
	}

	if l, err := zap.NewStdLogAt(logger.zapLogger, zap.ErrorLevel); err == nil {
		return l
	}

	return nil
}

func StdInfoLogger() *log.Logger {
	if logger == nil {
		return nil
	}

	if l, err := zap.NewStdLogAt(logger.zapLogger, zap.InfoLevel); err == nil {
		return l
	}

	return nil
}

func (z *zapLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	if checkedEntry := z.log.Check(zap.ErrorLevel, msg); checkedEntry != nil {
		checkedEntry.Write(handleFields(z.log, keysAndValues, zap.Error(err))...)
	}
}

func (z *zapLogger) Errorf(format string, args ...interface{}) {
	z.log.Sugar().Errorf(format, args...)
}

func V(level int) InfoLogger { return logger.V(level) }

func (z *zapLogger) V(level int) InfoLogger {
	if level < 0 || level > 1 {
		panic("valid log level is [0, 1]")
	}
	l := zapcore.Level(-1 * level)
	if z.zapLogger.Core().Enabled(l) {
		return &infoLogger{
			log:   z.zapLogger,
			level: l,
		}
	}

	return disabledInfoLogger
}

func (z *zapLogger) Write(p []byte) (n int, err error) {
	z.zapLogger.Info(string(p))

	return len(p), nil
}

func (z *zapLogger) WithValues(keysAndValues ...interface{}) Logger {
	newLogger := z.zapLogger.With(handleFields(z.zapLogger, keysAndValues)...)

	return NewLogger(newLogger)
}

func WithValues(keysAndValues ...interface{}) Logger {
	return logger.WithValues(keysAndValues...)
}

func (l *zapLogger) WithName(name string) Logger {
	newLogger := l.zapLogger.Named(name)

	return NewLogger(newLogger)
}

func WithName(name string) Logger {
	return logger.WithName(name)
}

func (z *zapLogger) Flush() {
	_ = z.zapLogger.Sync()
}

func Flush() {
	logger.Flush()
}

func ZapLogger() *zap.Logger {
	return logger.zapLogger
}

func NewLogger(l *zap.Logger) Logger {
	return &zapLogger{
		zapLogger: l,
		infoLogger: infoLogger{
			log:   l,
			level: zap.InfoLevel,
		},
	}
}

// message at the specified level is enabled.
func CheckIntLevel(level int32) bool {
	var lvl zapcore.Level
	if level < 5 {
		lvl = zapcore.InfoLevel
	} else {
		lvl = zapcore.DebugLevel
	}
	checkEntry := logger.zapLogger.Check(lvl, "")

	return checkEntry != nil
}

// Debug method output debug level log.
func Debug(msg string, fields ...Field) {
	logger.zapLogger.Debug(msg, fields...)
}

// Debugf method output debug level log.
func Debugf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Debugf(format, v...)
}

// Debugw method output debug level log.
func Debugw(msg string, keysAndValues ...interface{}) {
	logger.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

// Info method output info level log.
func Info(msg string, fields ...Field) {
	logger.zapLogger.Info(msg, fields...)
}

// Infof method output info level log.
func Infof(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Infof(format, v...)
}

// Infow method output info level log.
func Infow(msg string, keysAndValues ...interface{}) {
	logger.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

// Warn method output warning level log.
func Warn(msg string, fields ...Field) {
	logger.zapLogger.Warn(msg, fields...)
}

// Warnf method output warning level log.
func Warnf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Warnf(format, v...)
}

// Warnw method output warning level log.
func Warnw(msg string, keysAndValues ...interface{}) {
	logger.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

// Error method output error level log.
func Error(msg string, fields ...Field) {
	logger.zapLogger.Error(msg, fields...)
}

// Errorf method output error level log.
func Errorf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Errorf(format, v...)
}

// Errorw method output error level log.
func Errorw(msg string, keysAndValues ...interface{}) {
	logger.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

// Panic method output panic level log and shutdown application.
func Panic(msg string, fields ...Field) {
	logger.zapLogger.Panic(msg, fields...)
}

// Panicf method output panic level log and shutdown application.
func Panicf(format string, v ...interface{}) {
	logger.zapLogger.Sugar().Panicf(format, v...)
}

// Panicw method output panic level log.
func Panicw(msg string, keysAndValues ...interface{}) {
	logger.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

// Fatal method output fatal level log.
func Fatal(msg string, fields ...Field) {
	logger.zapLogger.Fatal(msg, fields...)
}

// Fatalf method output fatal level log.
func Fatalf(format string, args ...interface{}) {
	logger.zapLogger.Sugar().Fatalf(format, args...)
}

// Fatalw method output Fatalw level log.
func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.zapLogger.Fatal(msg, handleFields(logger.zapLogger, keysAndValues)...)
}

func GetOptions() *Options {
	return options
}

func GetLogger() *zapLogger {
	return logger
}
