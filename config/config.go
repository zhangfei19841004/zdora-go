package config

import (
	"github.com/koding/multiconfig"
	"sync"
	"zdora/constants"
)

type Config struct {
	Server ServerConfig
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}
		m := &multiconfig.DefaultLoader{}
		m.Loader = multiconfig.MultiLoader(newLoader(constants.ServiceName))
		m.Validator = multiconfig.MultiValidator(
			&multiconfig.RequiredValidator{},
		)
		err := m.Load(cfg)
		if err != nil {
			panic(err)
		}
	})
	return cfg
}

type ServerConfig struct {
	Host string `default:"localhost"`
	Port string `default:":9090"`
}
