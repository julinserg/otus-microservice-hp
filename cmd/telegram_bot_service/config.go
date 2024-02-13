package main

import "github.com/BurntSushi/toml"

type Config struct {
	Logger  LoggerConf
	AMQP    AMQPConfig
	TGBot   TGBotConfig
	AuthSrv AuthSrvConfig
	Debug   DebugConfig
}

type LoggerConf struct {
	Level string
}

type AMQPConfig struct {
	URI string
}

type TGBotConfig struct {
	Token   string
	Timeout int
}

type AuthSrvConfig struct {
	URI string
}

type DebugConfig struct {
	Host string
	Port string
}

func (c *Config) Read(fpath string) error {
	_, err := toml.DecodeFile(fpath, c)
	return err
}

func NewConfig() Config {
	return Config{}
}
