package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/koding/multiconfig"
	"github.com/spf13/viper"
	"os"
	"sync"
	"zdora/constants"
)

type Config struct {
	Server ServerConfig
}

var cfg *Config
var once sync.Once
var cfgSchema string

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GetConfigSchema() string {
	once.Do(func() {
		config := Config{}
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fmt.Println(path)
		if isExist(path + string(os.PathSeparator) + "zdora.yaml") {
			v := viper.New()
			v.SetConfigType("yaml")
			v.SetConfigName("zdora")
			v.AddConfigPath(path)
			err = v.ReadInConfig()
			if err != nil {
				panic(err)
			}
			err = v.Unmarshal(&config)
			if err != nil {
				panic(err)
			}
		}
		json, err := json.Marshal(&config)
		if err != nil {
			panic(err)
		}
		cfgSchema = string(json)
	})
	return cfgSchema
}

func SetConfig(config string) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	v := viper.New()
	v.SetConfigType("json")
	err = v.ReadConfig(bytes.NewReader([]byte(config)))
	if err != nil {
		panic(err)
	}
	v.SetConfigType("yaml")
	v.SetConfigName("zdora")
	v.AddConfigPath(path)
	err = v.SafeWriteConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
			err := v.WriteConfig()
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}

func loadYamlConfig() *Config {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	viper.AddConfigPath(path)
	viper.SetConfigName("zdora")
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig() // 搜索路径，并读取配置数据
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	cfg := Config{}
	viper.Unmarshal(&cfg)
	return &cfg
}

func loadEnvConfig() *Config {
	cfg := Config{}
	m := &multiconfig.DefaultLoader{}
	m.Loader = multiconfig.MultiLoader(newLoader(constants.ServiceName))
	m.Validator = multiconfig.MultiValidator(
		&multiconfig.RequiredValidator{},
	)
	err := m.Load(&cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}

func GetConfig() *Config {
	once.Do(func() {
		cfg = loadYamlConfig()
		if cfg == nil {
			cfg = loadEnvConfig()
		}
	})
	return cfg
}

type ServerConfig struct {
	Host string `default:"localhost"`
	Port string `default:"9090"`
}
