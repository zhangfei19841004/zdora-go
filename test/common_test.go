package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"testing"
	"zdora/client/util"
	"zdora/config"
)

func TestUUID(t *testing.T) {
	fmt.Println(uuid.Must(uuid.NewV4(), nil).String())
}

func TestIp(t *testing.T) {
	ip, _ := util.ExternalIP()
	fmt.Println(ip)
}

type StructDemo1 struct {
	Server3 string
	Server4 struct {
		Server5 string
	}
}

type StructDemo struct {
	Server struct {
		Server1 struct {
			Test1 string
		}
		Server2 string
	}
	Server0 StructDemo1
}

func TestStruct(t *testing.T) {
	s := StructDemo{}
	err := json.Unmarshal([]byte(`{"Server":{"Server1":{"Test1":"1"},"Server2":"2"},"Server0":{"Server3":"3","Server4":{"Server5":"5"}}}`), &s)
	if err != nil {
		panic(err)
	}
	v := viper.New()
	v.SetConfigType("json")
	err = v.ReadConfig(bytes.NewReader([]byte(`{"Server":{"Server1":{"Test1":"10"},"Server2":"20"},"Server0":{"Server3":"30","Server6":{"Server5":"50"}}}`)))
	if err != nil {
		panic(err)
	}
	v.SetConfigType("yaml")
	v.SetConfigName("zdora1")
	v.AddConfigPath(".")
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

func TestViper(t *testing.T) {
	s := config.GetConfigSchema()
	fmt.Println(s)
}
