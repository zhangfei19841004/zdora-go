package client

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"zdora/client/util"
	"zdora/config"
	"zdora/constants"
	"zdora/execmd"
	"zdora/model"
	zu "zdora/util"
)

var c *websocket.Conn
var done = make(chan struct{})

func read() {
	defer close(done)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		var msg model.Message
		zu.Unmarshal(string(message), &msg)
		switch msg.MessageType {
		case constants.COMMON:
			fmt.Printf(msg.Message)
		case constants.CLIENT_EXEC:
			handleExec(&msg)
		}
	}
}

func handleExec(msg *model.Message) {
	execMsg := msg.ExecMessage
	if execMsg.Combined {
		res := &model.Message{
			MessageType: constants.CLIENT_EXEC,
			ExecMessage: &model.ExecutorMessage{
				Command: execMsg.Command,
			},
		}
		pid, execReturnMsg := execmd.RunCmdCombined(execMsg.Command...)
		res.Message = execReturnMsg
		res.TargetClientId = msg.TargetClientId
		res.ExecMessage.Pid = pid
		res_ := zu.Marshal(res)
		Send(res_)
	} else {
		var logs = make(chan string, 1)
		ctx, cancel := context.WithCancel(context.Background())
		pid := execmd.RunCmd(ctx, cancel, logs, execMsg.Command...)
		Send(zu.Marshal(&model.Message{
			MessageType:    constants.BEGIN_EXEC,
			TargetClientId: msg.TargetClientId,
			ExecMessage: &model.ExecutorMessage{
				Command: execMsg.Command,
				Pid:     pid,
			},
		}))
		res := &model.Message{
			MessageType: constants.CLIENT_EXECING,
			ExecMessage: &model.ExecutorMessage{
				Command: execMsg.Command,
			},
		}
		go func(res *model.Message) {
			defer close(logs)
			for {
				select {
				case <-ctx.Done():
					Send(zu.Marshal(&model.Message{
						MessageType: constants.END_EXEC,
						ExecMessage: &model.ExecutorMessage{
							Pid: pid,
						},
					}))
					return
				case logMsg, ok := <-logs:
					if !ok {
						Send(zu.Marshal(&model.Message{
							MessageType: constants.END_EXEC,
							ExecMessage: &model.ExecutorMessage{
								Pid: pid,
							},
						}))
						return
					} else {
						res.Message = logMsg
						res.ExecMessage.Pid = pid
						Send(zu.Marshal(res))
					}
				}
			}
		}(res)
	}
}

func Send(message string) {
	err := c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func Run() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: config.GetConfig().Server.Host + ":" + config.GetConfig().Server.Port, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	var header = make(http.Header)
	clientId := util.GenerateClientId()
	log.Printf("client id is %s", clientId)
	header.Add("clientId", clientId)
	ip, err := util.ExternalIP()
	if err != nil {
		panic(err)
	}
	header.Add("ip", ip.String())
	header.Add("clientType", strconv.Itoa(constants.GO_CLIENT))
	c, _, err = websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	go read()
	/*ticker := time.NewTicker(time.Second)
	defer ticker.Stop()*/
	for {
		select {
		case <-done:
			return
		/*case t := <-ticker.C:
		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		if err != nil {
			log.Println("write:", err)
			return
		}*/
		case <-interrupt:
			log.Println("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			/*time.Sleep(time.Second)
			select {
			case <-done:
			case <-time.After(time.Second):
				fmt.Println("exit!!!!!")
			}*/
			return
		}
	}
}
