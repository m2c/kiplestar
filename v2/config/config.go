package config

import (
	"fmt"
	slog "github.com/m2c/kiplestar/commons/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	_ "path/filepath"
	"time"
)

/**
 * kipledb config
 */
type DataBaseConfig struct {
	Host        string        `yaml:"host"`
	User        string        `yaml:"user"`
	Pwd         string        `yaml:"password"`
	DbName      string        `yaml:"name"`
	DataBase    string        `yaml:"database"`
	Port        int           `yaml:"port"`
	MaxIdleCons int           `yaml:"maxIdleConns"`
	MaxOpenCons int           `yaml:"maxOpenConns"`
	MaxLifeTime time.Duration `yaml:"maxLifeTime"`
}
type kafkaConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type RedisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
	Name     string `yaml:"name"`
}
type NotifyServiceConfig struct {
	Url    string `yaml:"url"`
	Secret string `yaml:"secret"`
	AppKey string `yaml:"appKey"`
}

type Config struct {
	DataBase    []DataBaseConfig    `yaml:"dataSource"`
	Redis       []RedisConfig       `yaml:"redis"`
	Kafka       kafkaConfig         `yaml:"kafka"`
	Oss         OssConfig           `yaml:"oss"`
	Notify      NotifyServiceConfig `yaml:"notify"`
	RiskControl RiskControlConfig   `yaml:"riskControl"`
}

type OssConfig struct {
	OssBucket       string `yaml:"ossBucket"`
	AccessKeyID     string `yaml:"accessKeyID"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	OssEndPoint     string `yaml:"ossEndPoint"`
}

type RiskControlConfig struct {
	XApiKey    string `yaml:"xApiKey"`
	Host       string `yaml:"host"`
	Mock       bool   `yaml:"mock"`
	PortalHost string `yaml:"portalHost"`
}

const (
	EmailSendUrl = "/se/api/mail/sendmail"
)

func InitAllConfig(fileName string) *Config {
	var err error
	YamlFile, err = ioutil.ReadFile(fileName)
	if err != nil {
		slog.Info("load conf error, will exit")
		fmt.Println(err.Error())
		os.Exit(0)
	}
	dbc := &Config{}
	err = yaml.Unmarshal(YamlFile, dbc)
	if err != nil {
		fmt.Println(err.Error())
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
