package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	GetEnv `mapstructure:",squash"`
}

type GetEnv struct {
	Host    string `mapstructure:"TEST"`
	DbMysql string `mapstructure:"DB_CONNECTION_MYSQL"`
	DbPg    string `mapstructure:"DB_CONNECTION_PG"`
	Port    string `mapstructure:"PORT"`
}

func LoadConfig(fileName string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("app")
	v.AddConfigPath(fileName)
	v.SetConfigType("env")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		log.Println(err)
		return nil, errors.New("config not found")
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
	}
	return &c, err
}
