package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// log七种级别的日志
const (
	DebugLevel  = "Debug"
	InfoLevel   = "Info"
	WarnLevel   = "Warn"
	ErrorLevel  = "Error"
	DPanicLevel = "DPanic"
	PanicLevel  = "Panic"
	FatalLevel  = "Fatal"
)

var Log *zap.Logger
var sugaredLogger *zap.SugaredLogger
var initLogger *zap.SugaredLogger

var SSlogger *Slogger

var logger *Logger

type Logger struct {
	adapter *Adapter
}

type Adapter struct {
	//Path        string
	//Level       string
	//MaxFileSize int
	//MaxBackups  int
	//MaxAge      int
	//Compress    bool
	//Caller      bool

	suger *zap.SugaredLogger
}

type Slogger struct {
	suger *zap.SugaredLogger
}

type Loggor interface {
	Trace()
	Debug()
	Info()
	Warn()
	Error()
	/*
		Panic()
		Fatal()

		Tracef()
		Debugf()
		Infof()
		Warnf()
		Errorf()
		Fatalf()
		Panicf()
	*/
}

func InitLoggerFromParams(logFormat string, level string, logFile string, maxSize int, maxbackups int, maxAge int, localTime bool, compress bool) {

	/*
		logCore := setZapLoggerCore(logFormat, logFile, maxSize, maxAge, maxbackups, localTime, compress)

		Log := zap.New(logCore, zap.AddCaller())
		//logger, _ := zap.NewProduction()

		initLogger = Log.Sugar()
		//Logger = logger.Sugar()
		//initLogger = logger

		SSlogger = NewSlogger()
		SSlogger.suger = Log.Sugar()
		fmt.Println(initLogger)
		initLogger.Debugf("warning")
	*/
	adapLog := NewAdapter(logFormat, level, logFile, maxSize, maxbackups, maxAge, localTime, compress)
	initLogger = adapLog.suger
	logger = &Logger{
		adapter: adapLog,
	}
}

func NewAdapter(logFormat string, level string, logFile string, maxSize int, maxbackups int, maxAge int, localTime bool, compress bool) *Adapter {

	adplog := initLoggerFromParams(logFormat, level, logFile, maxSize, maxbackups, maxAge, localTime, compress)
	return &Adapter{
		suger: adplog,
	}
}

func initLoggerFromParams(logFormat string, level string, logFile string, maxSize int, maxbackups int, maxAge int, localTime bool, compress bool) *zap.SugaredLogger {

	logCore := setZapLoggerCore(logFormat, logFile, maxSize, maxAge, maxbackups, localTime, compress)
	Log := zap.New(logCore, zap.AddCaller(), zap.AddCallerSkip(3))
	//logger, _ := zap.NewProduction()

	initLogger := Log.Sugar()

	return initLogger

}

func NewSlogger() *Slogger {
	return &Slogger{}
}

func InitFromConfig(configType string, filePath string) {

}

// 定义日志打印的格式
func setZapLoggerEncoder(logFormat string) zapcore.Encoder {

	var encoder zapcore.Encoder

	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//encoderConfig.LevelEncoder = zapcore.CapitalColorLevelEncoder

	/*
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}

		encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}
	*/
	/*
		encoderConfig := zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			NameKey:      "mqtt",
			CallerKey:    "file",
			EncodeLevel:  zapcore.CapitalColorLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
			EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			}, //输出的时间格式
			EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendInt64(int64(d) / 1000000)
			},
		}
	*/

	if logFormat == "json" {

		encoder = zapcore.NewJSONEncoder(encoderConfig)

	} else {

		encoder = zapcore.NewConsoleEncoder(encoderConfig)

	}
	//encoder = zapcore.NewConsoleEncoder(encoderConfig)

	return encoder
}

// 定义日志的writer，这里我们使用lumback
func setZapLoggerLogFileWriter(logFile string, maxSize int, maxAge int, maxbackups int, localTime bool, compress bool) zapcore.WriteSyncer {

	logFileWriterConfig := &lumberjack.Logger{

		Filename:   logFile,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxbackups,
		LocalTime:  localTime,
		Compress:   compress,
	}

	logFileWriterSyncer := zapcore.AddSync(logFileWriterConfig)

	return logFileWriterSyncer

}

// 定义日志输出的级别
func setZapLoggerLevel(logLever string) (zapcore.Level, error) {

	l, err := zapcore.ParseLevel("debug")
	if err != nil {
		return 0, err
	}

	return l, nil

}

// 定义日志的core文件
func setZapLoggerCore(logFormat string, logFile string, maxSize int, maxAge int, maxbackups int, localTime bool, compress bool) zapcore.Core {

	encoder := setZapLoggerEncoder(logFormat)
	writeSyncer := setZapLoggerLogFileWriter(logFile, maxSize, maxAge, maxbackups, localTime, compress)
	//loggerLever :=
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	return core
}

// Debug 使用方法：log.Debug("test")
func (l *Logger) Debug(args ...interface{}) {

	//initLogger.Debug(args...)
	//Logger.Sugar().Debug(args...)
	//SSlogger.suger.Debug(args...)
	l.adapter.Debug(args...)

}

func (l *Logger) Info(args ...interface{}) {

	//initLogger.Info(args...)
	//Logger.Sugar().Info(args...)
	//SSlogger.suger.Debug(args...)
	l.adapter.Info(args...)

}

func (l *Logger) Warn(args ...interface{}) {

	//initLogger.Warn(args...)
	l.adapter.Warn(args...)

}

func (l *Logger) Error(args ...interface{}) {

	l.adapter.Error(args...)
}

func Panic(args ...interface{}) {

	initLogger.Panic(args...)

}

func Fatal(args ...interface{}) {
	initLogger.Fatal(args...)
}

// format output logs
func Debugf(template string, args ...interface{}) {

	initLogger.Debugf(template, args...)

}

func Infof(template string, args ...interface{}) {

	initLogger.Infof(template, args...)

}

func Warnf(template string, args ...interface{}) {

	initLogger.Warnf(template, args...)

}

func Errorf(template string, args ...interface{}) {

	initLogger.Errorf(template, args...)

}

func Panicf(template string, args ...interface{}) {

	initLogger.Panicf(template, args...)

}

func Fatalf(template string, args ...interface{}) {

	initLogger.Fatalf(template, args...)

}

func Sync() {

	//Logger.adapter.suger.sync()
}

///

func (adapter *Adapter) Debug(args ...interface{}) {
	adapter.suger.Debug(args...)
}

func (adapter *Adapter) Info(args ...interface{}) {
	adapter.suger.Info(args...)
}

func (adapter *Adapter) Warn(args ...interface{}) {
	adapter.suger.Warn(args...)
}

func (adapter *Adapter) Error(args ...interface{}) {
	adapter.suger.Error(args...)
}

func (adapter *Adapter) DPanic(args ...interface{}) {
	adapter.suger.DPanic(args...)
}

func (adapter *Adapter) Panic(args ...interface{}) {
	adapter.suger.Panic(args...)
}

func (adapter *Adapter) Fatal(args ...interface{}) {
	adapter.suger.Fatal(args...)
}

func (adapter *Adapter) Debugf(template string, args ...interface{}) {
	adapter.suger.Debugf(template, args...)
}

func (adapter *Adapter) Infof(template string, args ...interface{}) {
	adapter.suger.Infof(template, args...)
}

func (adapter *Adapter) Warnf(template string, args ...interface{}) {
	adapter.suger.Warnf(template, args...)
}

func (adapter *Adapter) Errorf(template string, args ...interface{}) {
	adapter.suger.Errorf(template, args...)
}

func (adapter *Adapter) DPanicf(template string, args ...interface{}) {
	adapter.suger.DPanicf(template, args...)
}

func (adapter *Adapter) Panicf(template string, args ...interface{}) {
	adapter.suger.Panicf(template, args...)
}

func (adapter *Adapter) Fatalf(template string, args ...interface{}) {
	adapter.suger.Fatalf(template, args...)
}

func (adapter *Adapter) Debugw(msg string, keysAndValues ...interface{}) {
	adapter.suger.Debugw(msg, keysAndValues...)
}

func (adapter *Adapter) Infow(msg string, keysAndValues ...interface{}) {
	adapter.suger.Infow(msg, keysAndValues...)
}

func (adapter *Adapter) Warnw(msg string, keysAndValues ...interface{}) {
	adapter.suger.Warnw(msg, keysAndValues...)
}

func (adapter *Adapter) Errorw(msg string, keysAndValues ...interface{}) {
	adapter.suger.Errorw(msg, keysAndValues...)
}

func (adapter *Adapter) DPanicw(msg string, keysAndValues ...interface{}) {
	adapter.suger.DPanicw(msg, keysAndValues...)
}

func (adapter *Adapter) Panicw(msg string, keysAndValues ...interface{}) {
	adapter.suger.Panicw(msg, keysAndValues...)
}

func (adapter *Adapter) Fatalw(msg string, keysAndValues ...interface{}) {
	adapter.suger.Fatalw(msg, keysAndValues...)
}
