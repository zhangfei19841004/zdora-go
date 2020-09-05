package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"zdora/client/util"
	"zdora/config"
	"zdora/constants"
	"zdora/model"
	zu "zdora/util"
)

func getCli() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: config.GetConfig().Server.Host + config.GetConfig().Server.Port, Path: "/ws"}
	var header = make(http.Header)
	clientId := util.GenerateClientId()
	header.Add("clientId", clientId)
	ip, err := util.ExternalIP()
	if err != nil {
		panic(err)
	}
	header.Add("ip", ip.String())
	header.Add("clientType", strconv.Itoa(constants.COBRA_CLIENT))
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func cliSendMsg(c *websocket.Conn, msg *model.Message) error {
	m := zu.Marshal(msg)
	err := c.WriteMessage(websocket.TextMessage, []byte(m))
	if err != nil {
		fmt.Printf("write close:%+v\n", err)
		return err
	}
	return nil
}

func Cli(msg *model.Message) ([]byte, error) {
	c, err := getCli()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	err = cliSendMsg(c, msg)
	if err != nil {
		return nil, err
	}
	_, message, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, err
}

func CliWithInterrupt(msg *model.Message, interrupt chan os.Signal) error {
	c, err := getCli()
	if err != nil {
		return err
	}
	defer c.Close()
	err = cliSendMsg(c, msg)
	if err != nil {
		return err
	}
	for {
		select {
		case <-interrupt:
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return nil
			}
		default:
			_, message, err := c.ReadMessage()
			if err != nil {
				return err
			}
			var msg model.Message
			zu.Unmarshal(string(message), &msg)
			switch msg.MessageType {
			case constants.CLOSE:
				fmt.Println(msg.Message)
				return nil
			default:
				fmt.Println(msg.Message)
			}

		}
	}
	return nil
}
