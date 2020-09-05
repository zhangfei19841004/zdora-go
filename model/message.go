package model

import (
	"zdora/types"
)

type Message struct {
	MessageType    int
	TargetClientId types.ClientId //目标clientId
	Message        string
	ExecMessage    *ExecutorMessage
}

type ExecutorMessage struct {
	Ip       types.IP //执行的clientip
	Command  []string
	Pid      types.PID
	Combined bool //是否联合结果输出，需要运行RunCmdCombined函数时用到
}

type ExecutorInfo struct {
	ClientId types.ClientId
	Pid      types.PID
	Command  []string
}

type PsCommandInfo struct {
	ClientId types.ClientId
	Ip       types.IP
	Pid      types.PID
	Commond  []string
}
