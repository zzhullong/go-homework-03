package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type AppConfig struct {
	DB DBConfig `yaml:"db"`
}

type DBConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

func Load() (*AppConfig, error) {
	// 加载本地.env文件
	_ = godotenv.Load()

	// 读取yaml
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// 开启系统环境变量，环境变量可以覆盖yaml
	viper.AutomaticEnv()
	viper.BindEnv("db.driver", "DB_DRIVER")
	viper.BindEnv("db.dsn", "DB_DSN")

	var cfg AppConfig
	err = viper.Unmarshal(&cfg)
	return &cfg, err
}
