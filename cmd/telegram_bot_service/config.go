package main

import "github.com/BurntSushi/toml"

type Config struct {
	Logger LoggerConf
	AMQP   AMQPConfig
	TGBot  TGBotConfig
	YDisk  YDiskConfig
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

type YDiskConfig struct {
	ClientId     string
	ClientSecret string
	Token        string
}

func (c *Config) Read(fpath string) error {
	_, err := toml.DecodeFile(fpath, c)
	return err
}

func NewConfig() Config {
	return Config{}
}
