package logging

import (
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level refers to the log logging level
type Level int8

// Create a general Base logger
var (
	baseLogger Logger
)

const (
	// Debug Level level. Usually only enabled when debugging. Very verbose logging.
	Debug Level = iota - 1

	// Info Level level. General operational entries about what's going on inside the
	// application.
	Info

	// Warn Level level. Non-critical entries that deserve eyes.
	Warn

	// Error Level level. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	Error

	// Panic Level level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	Panic

	// Fatal Level level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	Fatal
)

type WriteSyncer interface {
	io.Writer
	Sync() error
}

var once sync.Once

// Init needs to be called to ensure our logging has been initialized
func Init() {
	once.Do(func() {
		// By default, log to stderr (logrus's default), only warnings and above.
		baseLogger = NewLogger()
		baseLogger.SetLevel(Warn)
	})
}

func init() {
	Init()
}

// Fields maps logrus fields
type Field = interface{}

// Logger is the interface for loggers.
type Logger interface {
	// Debug logs a message at level Debug.
	Debug(...interface{})
	Debugln(...interface{})
	Debugf(string, ...interface{})

	// Info logs a message at level Info.
	Info(...interface{})
	Infoln(...interface{})
	Infof(string, ...interface{})

	// Warn logs a message at level Warn.
	Warn(...interface{})
	Warnln(...interface{})
	Warnf(string, ...interface{})

	// Error logs a message at level Error.
	Error(...interface{})
	Errorln(...interface{})
	Errorf(string, ...interface{})

	ErrorStack(args ...interface{})

	// Fatal logs a message at level Fatal.
	Fatal(...interface{})
	Fatalln(...interface{})
	Fatalf(string, ...interface{})

	// Panic logs a message at level Panic.
	Panic(...interface{})
	Panicln(...interface{})
	Panicf(string, ...interface{})

	// Add one key-value to log
	With(key string, value interface{}) Logger

	// WithFields logs a message with specific fields
	WithFields(...Field) Logger

	// Set the logging version (Info by default)
	SetLevel(Level)

	// Get the logging version
	GetLevel() Level

	// Set logger output
	SetOutput(w WriteSyncer)
}

type logger struct {
	log   *zap.SugaredLogger
	level *Level
	ws    *wrapWriter
}

func (l logger) With(key string, value interface{}) Logger {
	return &logger{
		log:   l.log.With(key, value),
		level: l.level,
		ws:    l.ws,
	}
}

func (l logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l logger) Debugln(args ...interface{}) {
	l.log.Debugln(args...)
}

func (l logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l logger) Infoln(args ...interface{}) {
	l.log.Infoln(args...)
}

func (l logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l logger) Warnln(args ...interface{}) {
	l.log.Warnln(args...)
}

func (l logger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l logger) ErrorStack(args ...interface{}) {
	l.log.With("stacktrace", string(debug.Stack())).Error(args...)
}

func (l logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l logger) Errorln(args ...interface{}) {
	l.log.Errorln(args...)
}

func (l logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l logger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l logger) Fatalln(args ...interface{}) {
	l.log.Fatalln(args...)
}

func (l logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l logger) Panic(args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	l.log.Panic(args...)
}

func (l logger) Panicln(args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	l.log.Panicln(args...)
}

func (l logger) Panicf(format string, args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	l.log.Panicf(format, args...)
}

func (l logger) WithFields(fields ...Field) Logger {
	return &logger{
		log:   l.log.With(fields),
		level: l.level,
		ws:    nil,
	}
}

func (l logger) GetLevel() (lvl Level) {
	return Level(l.log.Level())
}

func (l *logger) SetLevel(lvl Level) {
	l.level = &lvl
}

func (l logger) SetOutput(w WriteSyncer) {
	l.ws.writerSyncer = w
}

// Base returns the default Logger logging to
func Base() Logger {
	return baseLogger
}

// NewLogger returns a new Logger logging to out.
func NewLogger() Logger {
	level := Warn
	l := &logger{
		level: &level,
		ws: &wrapWriter{
			writerSyncer: os.Stdout,
		},
		log: nil,
	}

	levelFunc := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level <= zapcore.Level(*l.level)
	})

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, l.ws, levelFunc)

	zlog := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.log = zlog.Sugar()

	return l
}

type wrapWriter struct {
	writerSyncer WriteSyncer
}

func (w *wrapWriter) Write(bs []byte) (n int, err error) {
	return w.writerSyncer.Write(bs)
}

func (w *wrapWriter) Sync() error {
	return w.writerSyncer.Sync()
}

// RegisterExitHandler registers a function to be called on exit by logrus
// Exit handling happens when logrus.Exit is called, which is called by logrus.Fatal
func RegisterExitHandler(handler func()) {
	// TODO
	logrus.RegisterExitHandler(handler)
}
