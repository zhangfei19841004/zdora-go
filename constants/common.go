package constants

var (
	ServiceName = "zdora"
)

const (
	GO_CLIENT = 1 + iota
	COBRA_CLIENT
)

//message type
const (
	PS             = 10001 + iota //ps命令
	COMMON                        //普通消息
	ZDORA_EXEC                    //zdora发起命令
	CLIENT_EXEC                   //client执行命令
	CLIENT_EXECING                //client执行命令
	LOGS                          //client执行log信息
	CLOSE                         //服务端发消息给客户端，让客户端关闭
	BEGIN_EXEC
	END_EXEC
)
