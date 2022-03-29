package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"sync"
)

var loggerChannel = make(map[string]*Logger)
var lock sync.Mutex

type Logger struct {
	zap        *zap.Logger
	channel    string
	filename   string
	filepath   string
	tempParams []string
}

type ConfigOption struct {
	dingEnable   bool   `json:"ding_enable,omitempty" yaml:"ding_enable"`     // 钉钉通知
	TraceEnable  bool   `json:"trace_enable,omitempty" yaml:"trace_enable"`   // 日志跟踪
	ConsolePrint bool   `json:"console_print,omitempty" yaml:"console_print"` // 控制台显示
	Level        string `json:"level,omitempty" yaml:"level"`
	LogPath      string `json:"log_path,omitempty" yaml:"log_path"`
	MaxSize      int    `json:"max_size,omitempty" yaml:"max_size"`
	MaxBackups   int    `json:"max_backups,omitempty" yaml:"max_backups"`
	MaxAge       int    `json:"max_age,omitempty" yaml:"max_age"`
	Compress     bool   `json:"compress,omitempty" yaml:"compress"`
}

var config = ConfigOption{
	TraceEnable:  true,
	ConsolePrint: true,
	Level:        "debug",
	LogPath:      "./log",
	MaxSize:      1024,
	MaxBackups:   10,
	MaxAge:       0,
	Compress:     false,
}

func InitConfig(c *ConfigOption) {
	if c == nil {
		return
	}
	config.dingEnable = c.dingEnable
	config.TraceEnable = c.TraceEnable
	config.ConsolePrint = c.ConsolePrint
	config.Compress = c.Compress

	if c.LogPath != "" {
		config.LogPath = c.LogPath
	}
	if c.Level != "" {
		config.Level = c.Level
	}
	if c.MaxSize > 0 {
		config.MaxSize = c.MaxSize
	}
	if c.MaxAge > 0 {
		config.MaxAge = c.MaxAge
	}
	if c.MaxBackups > 0 {
		config.MaxBackups = c.MaxBackups
	}
}

// Channel 日志通道
func Channel(filename string, dirs ...string) *Logger {
	return getLogger(filename, dirs...)
}

func getLogger(channel string, subDir ...string) *Logger {
	if channel == "" {
		channel = "server" //默认日志文件
	}

	key := channel + strings.Join(subDir, "_")
	if loggerChannel[key] != nil {
		return loggerChannel[key]
	}

	lock.Lock()
	defer lock.Unlock()

	//初始化
	this := &Logger{
		channel:  channel,
		filename: strings.Trim(channel, "/") + ".log",
		filepath: strings.Join(subDir, "/"), //日志目录下的子目录,默认channel
	}

	loggerChannel[key] = this.initZapLogger()
	return this
}

// Info 提示级别
func Info(msg string, fields ...zap.Field) *Logger {
	return getLogger("info").Info(msg, fields...)
}
func (l *Logger) Info(msg string, fields ...zap.Field) *Logger {
	defer l.zap.Sync()
	l.zap.Info(msg, fields...)
	return l
}

// Error 错误级别
func Error(msg string, fields ...zap.Field) *Logger {
	return getLogger("error").Error(msg, fields...)
}
func (l *Logger) Error(msg string, fields ...zap.Field) *Logger {
	defer l.zap.Sync()
	l.tempParams = []string{msg}
	l.zap.Error(msg, fields...)
	return l
}

// Debug 调试级别
func Debug(msg string, fields ...zap.Field) *Logger {
	return getLogger("debug").Debug(msg, fields...)
}
func (l *Logger) Debug(msg string, fields ...zap.Field) *Logger {
	defer l.zap.Sync()
	l.tempParams = []string{msg}
	l.zap.Debug(msg, fields...)
	return l
}

// 初始化 zap log
func (l *Logger) initZapLogger() *Logger {
	// 日志文件配置

	hook := lumberjack.Logger{
		Filename:   l.GetFilepath(),   // 日志文件路径
		MaxSize:    config.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: config.MaxBackups, // 日志文件最多保存多少个备份
		MaxAge:     config.MaxAge,     // 文件最多保存多少天
		Compress:   config.Compress,   // 是否压缩
	}
	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch strings.ToLower(config.Level) {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	case "panic":
		atomicLevel.SetLevel(zap.PanicLevel)
	default:
		atomicLevel.SetLevel(zap.DebugLevel)
	}
	/*** 配置日志 ***/
	var logOptions []zap.Option

	// 设置日志级别
	logWriteSyncer := []zapcore.WriteSyncer{
		zapcore.AddSync(&hook), //默认输出到文件
	}

	// 调试模式 打印到控制台
	if config.ConsolePrint == true {
		logWriteSyncer = append(logWriteSyncer, zapcore.AddSync(os.Stdout))
	}

	// 调试模式，堆栈跟踪
	if config.TraceEnable == true {
		logOptions = append(logOptions, zap.AddCaller())
		logOptions = append(logOptions, zap.Development())
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(logWriteSyncer...),
		atomicLevel,
	)

	// 构造日志
	l.zap = zap.New(core, logOptions...)
	return l
}

func (l *Logger) GetFilepath() string {
	//日志默认配置目录
	filename := strings.TrimRight(config.LogPath, "/")
	if filename != "" {
		filename += "/"
	}
	//拼接用户目录
	if strings.Trim(l.filepath, " ") != "" {
		filename += l.filepath + "/"
	}
	//拼接文件名
	filename += l.filename
	return filename
}
