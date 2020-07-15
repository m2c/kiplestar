package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	slog "kiple_star/commons/log"
	"os"
	_ "path/filepath"
	"time"
)

/**
 * kipledb config
 */
type dataBaseConfig struct {
	Host        string        `yaml:"host"`
	User        string        `yaml:"user"`
	Pwd         string        `yaml:"password"`
	DbName      string        `yaml:"database"`
	Port        int           `yaml:"port"`
	MaxIdleCons int           `yaml:"maxIdleConns"`
	MaxOpenCons int           `yaml:"maxOpenConns"`
	MaxLifeTime time.Duration `yaml:"maxLifeTime"`
}
type kafkaConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type redisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}
type Config struct {
	DataBase dataBaseConfig `yaml:"dataSource"`
	Redis    redisConfig    `yaml:"qr_redis"`
	Kafka    kafkaConfig    `yaml:"kafka"`
}

func InitAllConfig(fileName string) *Config {
	var err error
	YamlFile, err = ioutil.ReadFile(fileName)
	if err != nil {
		slog.Info("load conf error, will exit")
	}
	dbc := &Config{}
	err = yaml.Unmarshal(YamlFile, dbc)
	if err != nil {
		slog.Info(err.Error())
		os.Exit(0)
	}
	return dbc
}

func LoadCustomizeConfig(config interface{}) error {
	err := yaml.Unmarshal(YamlFile, config)
	if err != nil {
		return err
	}
	return nil
}
