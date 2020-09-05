package util

import (
	"encoding/json"
	"fmt"
)

func Marshal(msg interface{}) string {
	bs, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("err:%+v\n", err)
		return ""
	}
	return string(bs)
}

func Unmarshal(msg string, message interface{}) {
	err := json.Unmarshal([]byte(msg), message)
	if err != nil {
		fmt.Printf("err:%+v\n", err)
	}
}
