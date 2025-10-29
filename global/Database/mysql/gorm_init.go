// FilePath: C:/WoopsBBS/global/Database/mysql\gorm_init.go
package mysql

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"

	"WoopsBBS/global/model"
)

type dbConfig struct {
	Username        string          `mapstructure:"username"`
	Password        string          `mapstructure:"password"`
	Host            string          `mapstructure:"host"`
	Port            string          `mapstructure:"port"`
	DBName          string          `mapstructure:"dbname"`
	Charset         string          `mapstructure:"charset"`
	MaxOpenConns    int             `mapstructure:"max_open_conns"`
	MaxIdleConns    int             `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration   `mapstructure:"conn_max_lifetime"` // 单位：秒
	LogMode         logger.LogLevel `mapstructure:"log_mode"`
}

// 全局DB实例（实际项目中可根据需求调整作用域）
var DB *gorm.DB

func GormInit() bool {
	// 1. 读取配置文件
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("配置文件加载失败: %v\n", err)
		return false
	}

	// 2. 初始化数据库连接
	db, err := initDB(cfg)
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		return false
	}

	// 3. 赋值全局DB实例
	DB = db

	// 4. 执行数据库迁移
	if err := AutoMigrate(DB); err != nil {
		fmt.Printf("数据库迁移失败: %v\n", err)
		return false
	}

	fmt.Println("GORM初始化成功")
	return true
}

// AutoMigrate 自动迁移数据库模型
func AutoMigrate(db *gorm.DB) error {
	fmt.Println("开始数据库迁移...")
	
	// 迁移用户模型
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		return fmt.Errorf("迁移User模型失败: %w", err)
	}
	
	fmt.Println("数据库迁移完成")
	return nil
}

// 加载配置文件（内部函数）
func loadConfig() (*dbConfig, error) {
	// 配置Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}

	// 解析配置到结构体
	var cfg dbConfig
	if err := viper.UnmarshalKey("database", &cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 转换时间单位（秒→纳秒）
	cfg.ConnMaxLifetime *= time.Second

	return &cfg, nil
}

// 初始化数据库连接（内部函数）
func initDB(cfg *dbConfig) (*gorm.DB, error) {
	// 校验必要参数
	if cfg.Username == "" || cfg.Password == "" || cfg.Host == "" || cfg.Port == "" || cfg.DBName == "" {
		return nil, errors.New("缺少必要的数据库配置参数")
	}

	// 默认字符集
	if cfg.Charset == "" {
		cfg.Charset = "utf8mb4"
	}

	// 构建DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	// 日志级别处理
	logLevel := cfg.LogMode
	if logLevel < logger.Silent || logLevel > logger.Info {
		logLevel = logger.Warn // 默认警告级别
	}

	// GORM配置
	gormCfg := &gorm.Config{
		Logger:                 logger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层连接失败: %w", err)
	}

	// 连接池参数
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	} else {
		sqlDB.SetMaxOpenConns(100)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	} else {
		sqlDB.SetMaxIdleConns(20)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	} else {
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("测试连接失败: %w", err)
	}

	return db, nil
}
