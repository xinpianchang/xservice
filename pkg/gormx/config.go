package gormx

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	_ "gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	gormopentracing "gorm.io/plugin/opentracing"

	"github.com/xinpianchang/xservice/pkg/log"
)

var (
	dbs map[string]*gorm.DB
)

type DbConfig struct {
	Name                         string `yaml:"name"`
	Uri                          string `yaml:"uri"`
	MaxConn                      int    `yaml:"maxConn"`
	MaxIdleConn                  int    `yaml:"maxIdleConn"`
	ConnMaxLifetimeInMillisecond int    `yaml:"connMaxLifetimeInMillisecond"`
	QueryFields                  bool   `yaml:"queryFields"`
	CreateBatchSize              int    `yaml:"createBatchSize"`
}

type ConfigureFn func(DbConfig) *gorm.DB

// Config config db, default use mysql
func Config(v *viper.Viper, configureFn ...ConfigureFn) {
	var cfg []DbConfig
	if err := v.UnmarshalKey("database", &cfg); err != nil {
		log.Fatal("read database config", zap.Error(err))
	}

	dbs = make(map[string]*gorm.DB, len(cfg))
	for _, c := range cfg {
		if c.MaxConn <= 0 {
			c.MaxConn = 100
		}

		if c.MaxIdleConn <= 0 {
			c.MaxIdleConn = 0
		}

		if c.ConnMaxLifetimeInMillisecond <= 0 {
			c.ConnMaxLifetimeInMillisecond = int((time.Minute * 5).Milliseconds())
		}

		if c.CreateBatchSize <= 0 {
			c.CreateBatchSize = 1000
		}

		var db *gorm.DB
		if len(configureFn) > 0 && configureFn[0] != nil {
			db = configureFn[0](c)
		} else {
			db = MySQLDbConfig(c)
		}
		if db == nil {
			continue
		}
		if err := db.Use(gormopentracing.New()); err != nil {
			log.Error("apply db opentracing", zap.Error(err))
		}
		dbs[c.Name] = db
	}
}

// MySQLDbConfig for mysql config
func MySQLDbConfig(cfg DbConfig) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.Uri), &gorm.Config{
		PrepareStmt:     true,
		QueryFields:     cfg.QueryFields,
		CreateBatchSize: cfg.CreateBatchSize,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("open db failed", zap.String("name", cfg.Name), zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("get db failed", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Millisecond * time.Duration(cfg.ConnMaxLifetimeInMillisecond))

	logger, _ := log.NewLogger(fmt.Sprint("sql-", cfg.Name, ".log"))
	db.Logger = &dbLogger{logger: logger.Named(cfg.Name)}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("db ping failed", zap.Error(err))
	}

	return db
}

func Get(name string) *gorm.DB {
	return dbs[name]
}
