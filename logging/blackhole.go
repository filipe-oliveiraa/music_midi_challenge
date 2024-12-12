package logging

type BlackHoleLogger struct{}

// ErrorStack implements Logger.
func (b *BlackHoleLogger) ErrorStack(args ...interface{}) {
	panic("unimplemented")
}

func NewBlackholeLogger() Logger {
	return &BlackHoleLogger{}
}

// Debug implements logging.Logger.
func (b *BlackHoleLogger) Debug(...interface{}) {
}

// Debugf implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Debugf(string, ...interface{}) {
}

// Debugln implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Debugln(...interface{}) {
}

// Error implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Error(...interface{}) {
}

// Errorf implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Errorf(string, ...interface{}) {
}

// Errorln implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Errorln(...interface{}) {
}

// Fatal implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Fatal(...interface{}) {
}

// Fatalf implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Fatalf(string, ...interface{}) {
}

// Fatalln implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Fatalln(...interface{}) {
}

// GetLevel implements logging.Logger that does nothing.
func (b *BlackHoleLogger) GetLevel() Level {
	return Info
}

// Info implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Info(...interface{}) {
}

// Infof implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Infof(string, ...interface{}) {
}

// Infoln implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Infoln(...interface{}) {
}

// IsLevelEnabled implements logging.Logger that does nothing.
func (b *BlackHoleLogger) IsLevelEnabled(level Level) bool {
	return false
}

// Panic implements logging.Loggerthat does nothing.
func (b *BlackHoleLogger) Panic(...interface{}) {
}

// Panicf implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Panicf(string, ...interface{}) {
}

// Panicln implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Panicln(...interface{}) {
}

// SetLevel implements logging.Logger that does nothing.
func (b *BlackHoleLogger) SetLevel(Level) {
}

// SetOutput implements logging.Logger that does nothing.
func (b *BlackHoleLogger) SetOutput(WriteSyncer) {
}

// Warn implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Warn(...interface{}) {
}

// Warnf implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Warnf(string, ...interface{}) {
}

// Warnln implements logging.Logger that does nothing.
func (b *BlackHoleLogger) Warnln(...interface{}) {
}

// With implements logging.Logger that does nothing.
func (b *BlackHoleLogger) With(key string, value interface{}) Logger {
	return b
}

// WithFields implements logging.Logger that does nothing.
func (b *BlackHoleLogger) WithFields(f ...Field) Logger {
	return b
}
