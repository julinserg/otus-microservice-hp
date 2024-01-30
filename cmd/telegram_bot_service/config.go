package main

import "github.com/BurntSushi/toml"

type Config struct {
	Logger LoggerConf
	AMQP   AMQPConfig
	TGBot  TGBotConfig
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

func (c *Config) Read(fpath string) error {
	_, err := toml.DecodeFile(fpath, c)
	return err
}

func NewConfig() Config {
	return Config{}
}
