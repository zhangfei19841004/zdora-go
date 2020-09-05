package test

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
	"zdora/client/util"
)

func TestUUID(t *testing.T) {
	fmt.Println(uuid.Must(uuid.NewV4(), nil).String())
}

func TestIp(t *testing.T) {
	ip, _ := util.ExternalIP()
	fmt.Println(ip)
}
