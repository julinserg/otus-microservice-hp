package main

import "github.com/BurntSushi/toml"

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	AMQP    AMQPConfig
	Debug   DebugConfig
	AuthSrv AuthSrvConfig
	Storage CloudStorageConfig
}

type LoggerConf struct {
	Level string
	// TODO
}

type AMQPConfig struct {
	URI string
}

type DebugConfig struct {
	TokenYD string
}

type AuthSrvConfig struct {
	URI string
}

type CloudStorageConfig struct {
	Folder string
}

func (c *Config) Read(fpath string) error {
	_, err := toml.DecodeFile(fpath, c)
	return err
}

func NewConfig() Config {
	return Config{}
}
