package zap

// Logger is a stub for go.uber.org/zap.Logger.
type Logger struct{}

// Field is a stub for zap.Field.
type Field struct{}

func NewProduction() (*Logger, error)  { return &Logger{}, nil }
func NewDevelopment() (*Logger, error) { return &Logger{}, nil }

func (l *Logger) Info(msg string, fields ...Field)   {}
func (l *Logger) Warn(msg string, fields ...Field)   {}
func (l *Logger) Error(msg string, fields ...Field)  {}
func (l *Logger) Debug(msg string, fields ...Field)  {}
func (l *Logger) DPanic(msg string, fields ...Field) {}
func (l *Logger) Panic(msg string, fields ...Field)  {}
func (l *Logger) Fatal(msg string, fields ...Field)  {}

func (l *Logger) Infof(template string, args ...interface{})   {}
func (l *Logger) Warnf(template string, args ...interface{})   {}
func (l *Logger) Errorf(template string, args ...interface{})  {}
func (l *Logger) Debugf(template string, args ...interface{})  {}
func (l *Logger) DPanicf(template string, args ...interface{}) {}
func (l *Logger) Panicf(template string, args ...interface{})  {}
func (l *Logger) Fatalf(template string, args ...interface{})  {}

func (l *Logger) Infow(msg string, keysAndValues ...interface{})   {}
func (l *Logger) Warnw(msg string, keysAndValues ...interface{})   {}
func (l *Logger) Errorw(msg string, keysAndValues ...interface{})  {}
func (l *Logger) Debugw(msg string, keysAndValues ...interface{})  {}
func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {}
func (l *Logger) Panicw(msg string, keysAndValues ...interface{})  {}
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{})  {}

func (l *Logger) Sugar() *SugaredLogger { return &SugaredLogger{} }

// SugaredLogger is a stub for zap.SugaredLogger.
type SugaredLogger struct{}

func (s *SugaredLogger) Info(args ...interface{})  {}
func (s *SugaredLogger) Warn(args ...interface{})  {}
func (s *SugaredLogger) Error(args ...interface{}) {}
func (s *SugaredLogger) Debug(args ...interface{}) {}

func (s *SugaredLogger) Infof(template string, args ...interface{})  {}
func (s *SugaredLogger) Warnf(template string, args ...interface{})  {}
func (s *SugaredLogger) Errorf(template string, args ...interface{}) {}
func (s *SugaredLogger) Debugf(template string, args ...interface{}) {}

func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{})  {}
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...interface{})  {}
func (s *SugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {}
func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{}) {}
