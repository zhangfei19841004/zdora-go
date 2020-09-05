package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strconv"
	"zdora/config"
	"zdora/constants"
	"zdora/model"
	"zdora/types"
	"zdora/util"
)

type wsServer struct {
	listener net.Listener
	addr     string
	upgrade  *websocket.Upgrader
}

func NewWsServer() *wsServer {
	ws := new(wsServer)
	ws.addr = config.GetConfig().Server.Port
	ws.upgrade = &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if r.Method != "GET" {
				fmt.Println("method is not GET")
				return false
			}
			if r.URL.Path != "/ws" {
				fmt.Println("path error")
				return false
			}
			return true
		},
	}
	return ws
}

func (ws *wsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ws" {
		httpCode := http.StatusInternalServerError
		reasePhrase := http.StatusText(httpCode)
		fmt.Println("path error ", reasePhrase)
		http.Error(w, reasePhrase, httpCode)
		return
	}
	conn, err := ws.upgrade.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("websocket error:", err)
		return
	}
	ws.listenCilents(r, conn)
}

func (ws *wsServer) listenCilents(r *http.Request, conn *websocket.Conn) {
	header := r.Header
	clientId := types.ClientId(header.Get("clientId"))
	ip := types.IP(header.Get("ip"))
	clientType, err := strconv.Atoi(header.Get("clientType"))
	if err != nil {
		fmt.Printf("error:%+v\n", err)
		conn.Close()
		return
	}
	info := &WsInfo{
		ClientId:   clientId,
		Ip:         ip,
		ClientType: clientType,
		Conn:       conn,
		StopCh:     make(chan int),
	}
	clients.Add(clientId, info)
	welcome := fmt.Sprintf("welcome [%s] [%s]", info.ClientId, info.Ip)
	fmt.Println(welcome)
	if !info.IsCobraClient() {
		msg := &model.Message{
			MessageType: constants.COMMON,
			Message:     welcome,
		}
		info.send(util.Marshal(msg))
	}
	go info.onMessage()
}

func Send(clientId, message string) {
	info, err := GetWsInfo(types.ClientId(clientId))
	if err != nil {
		fmt.Printf("message[%s]没有发送 - %s - [%+v]\n", message, clientId, err)
		return
	}
	info.send(message)
}

func (ws *wsServer) Start() (err error) {
	ws.listener, err = net.Listen("tcp", ws.addr)
	if err != nil {
		fmt.Println("net listen error:", err)
		return
	}

	err = http.Serve(ws.listener, ws)
	if err != nil {
		fmt.Println("http serve error:", err)
		return
	}
	return nil
}
