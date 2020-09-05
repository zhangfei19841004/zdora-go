package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"time"
	"zdora/constants"
	"zdora/model"
	"zdora/types"
	"zdora/util"
)

type WsInfo struct {
	ClientId   types.ClientId
	Ip         types.IP
	ClientType int
	Conn       *websocket.Conn `json:"-"`
	StopCh     chan int        `json:"-"`
}

func (info *WsInfo) onMessage() {
	defer func() {
		clients.Delete(info.ClientId)
		if info.IsCobraClient() {
			logsClis.Delete(info.ClientId)
		} else {
			execClis.DeleteInfo(info.ClientId, 0)
			logsClis.DeleteInfo(info.ClientId, 0)
		}
		info.Conn.Close()
	}()
	for {
		//info.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(5000)))
		_, msg, err := info.Conn.ReadMessage()
		if err != nil {
			close(info.StopCh)
			// 判断是不是超时
			if netErr, ok := err.(net.Error); ok {
				if netErr.Timeout() {
					fmt.Printf("ReadMessage timeout remote: %v\n", info.Conn.RemoteAddr())
					return
				}
			}
			if websocket.IsCloseError(err, websocket.CloseMessage, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				if !info.IsCobraClient() {
					fmt.Printf("%s已经下线!\n", info.ClientId)
				}
				return
			}
			// 其他错误，如果是 1001 和 1000 就不打印日志
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				if !info.IsCobraClient() {
					fmt.Printf("%s已经下线!\n", info.ClientId)
				}
				return
			}
			return
		}
		message := &model.Message{}
		util.Unmarshal(string(msg), message)
		switch message.MessageType {
		case constants.PS:
			handlePs(info)
		case constants.ZDORA_EXEC:
			handleExec(info, message)
		case constants.CLIENT_EXEC:
			handleClientExecing(info, message, 0)
		case constants.CLIENT_EXECING:
			handleClientExecing(info, message, 1)
		case constants.BEGIN_EXEC:
			handleClientBeginExec(info, message)
		case constants.END_EXEC:
			handleClientEndExec(info, message)
		case constants.LOGS:
			handleLogs(info, message)
		}

	}
}

func (info *WsInfo) monitor() {
	for {
		select {
		case <-info.StopCh:
			fmt.Println("connect closed")
			clients.Delete(info.ClientId)
			return
		case <-time.After(time.Second * 30): //心跳
			data := fmt.Sprintf("hello websocket test from server %v", time.Now().UnixNano())
			err := info.Conn.WriteMessage(1, []byte(data))
			fmt.Println("sending....")
			if err != nil {
				fmt.Println("send msg faild ", err)
				return
			}
		}
	}
}

func (info *WsInfo) send(message string) {
	err := info.Conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		fmt.Printf("send msg[%s] faild [%+v] \n", message, err)
		return
	}
}

func (info *WsInfo) IsCobraClient() bool {
	if info.ClientType == constants.COBRA_CLIENT {
		return true
	}
	return false
}
