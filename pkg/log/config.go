package log

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/xinpianchang/xservice/pkg/signalx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger    Logger
	zaplogger *zap.Logger
	cfg       Cfg

	loggerFileMap      map[string]*zap.Logger = make(map[string]*zap.Logger)
	loggerFileMapMutex sync.Mutex
)

// Cfg is the log config
type Cfg struct {
	Level   string `yaml:"level"`
	File    string `yaml:"file"`
	Stdout  bool   `yaml:"stdout"`
	Format  string `yaml:"format"`
	Caller  bool   `yaml:"caller"`
	MaxSize int    `yaml:"maxSize"`
	MaxDays int    `yaml:"maxDays"`
}

func init() {
	v := viper.New()
	v.SetDefault("log", map[string]interface{}{
		"level":  "info",
		"format": "console",
		"stdout": true,
		"caller": true,
	})
	Config(v)
}

// Config for log configuration
func Config(v *viper.Viper) {
	if err := v.UnmarshalKey("log", &cfg); err != nil {
		Fatal("parse log config", zap.Error(err))
	}

	if cfg.Level == "" {
		cfg.Level = "info"
	}

	if cfg.Format == "" {
		cfg.Format = "console"
	}

	if l, err := buildZapLogger(cfg); err != nil {
		Fatal("config log", zap.Error(err))
	} else {
		zaplogger = l
	}

	logger = newLogger(zaplogger)
}

// NewLoggerFile create a new logger with file and use the global log configuration
func NewLogger(file string) (Logger, error) {
	c := cfg
	if c.File == file {
		return logger, nil
	}

	if file != "" && c.File != "" {
		c.File = filepath.Join(filepath.Dir(cfg.File), file)
	}

	if l, err := buildZapLogger(cfg); err != nil {
		return nil, err
	} else {
		return newLogger(l), nil
	}
}

func buildZapLogger(cfg Cfg) (*zap.Logger, error) {
	loggerFileMapMutex.Lock()
	defer loggerFileMapMutex.Unlock()

	if l, ok := loggerFileMap[cfg.File]; ok {
		return l, nil
	}

	if cfg.Level == "" {
		cfg.Level = "info"
	}

	if cfg.Format == "" {
		cfg.Format = "console"
	}

	ws := make([]zapcore.WriteSyncer, 0, 2)

	if cfg.Stdout && os.Getenv("XSERVICE_DISABLE_STDOUT") == "" {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}

	if cfg.File != "" {
		if cfg.MaxSize <= 0 {
			cfg.MaxSize = 1024
		}

		if cfg.MaxDays <= 0 {
			cfg.MaxDays = 7
		}
	}

	if cfg.File != "" {
		rotateLogger := &lumberjack.Logger{
			Filename:  cfg.File,
			MaxSize:   cfg.MaxSize,
			MaxAge:    cfg.MaxDays,
			LocalTime: true,
			Compress:  true,
		}
		ws = append(ws, zapcore.AddSync(rotateLogger))

		go scheduleRotate(rotateLogger)
	}

	var level zapcore.Level
	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		Error("parse level", zap.Error(err))
		level = zapcore.InfoLevel
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	writeSynced := zapcore.NewMultiWriteSyncer(ws...)

	encoding := zap.NewProductionEncoderConfig()
	encoding.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if strings.ToLower(cfg.Format) == "json" {
		encoder = zapcore.NewJSONEncoder(encoding)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoding)
	}
	core := zapcore.NewCore(encoder, writeSynced, atomicLevel)

	options := make([]zap.Option, 0, 3)
	options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	if cfg.Caller {
		options = append(options, zap.AddCaller(), zap.AddCallerSkip(2))
	}
	log := zap.New(core, options...)
	if cfg.File != "" {
		signalx.AddShutdownHook(func(os.Signal) {
			log.Sync()
		})
	}
	return log, nil
}

func scheduleRotate(log *lumberjack.Logger) {
	for {
		n := time.Now().Add(time.Hour * 24)
		next := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.Local)
		d := time.Until(next)
		time.Sleep(d)
		_ = log.Rotate()
	}
}
