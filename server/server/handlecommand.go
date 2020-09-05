package server

import (
	"fmt"
	"zdora/constants"
	"zdora/model"
	"zdora/types"
	"zdora/util"
)

func handlePs(info *WsInfo) {
	infos := ListAllClients()
	infoString := util.Marshal(infos)
	info.send(infoString)
}

func handleExec(info *WsInfo, msg *model.Message) {
	fmt.Printf("%s\n", util.Marshal(msg))
	if msg.TargetClientId != "" {
		tinfo, err := GetWsInfo(msg.TargetClientId)
		if err != nil {
			info.send(err.Error())
			return
		}
		tmsg := &model.Message{
			MessageType:    constants.CLIENT_EXEC,
			TargetClientId: info.ClientId,
			ExecMessage: &model.ExecutorMessage{
				Command:  msg.ExecMessage.Command,
				Combined: msg.ExecMessage.Combined,
			},
		}
		tm := util.Marshal(tmsg)
		tinfo.send(tm)
	}
}

func handleClientExecing(info *WsInfo, msg *model.Message, execType int32) {
	tmsg := &model.Message{
		MessageType:    constants.COMMON,
		Message:        msg.Message,
		TargetClientId: info.ClientId,
		ExecMessage: &model.ExecutorMessage{
			Ip:      info.Ip,
			Command: msg.ExecMessage.Command,
			Pid:     msg.ExecMessage.Pid,
		},
	}
	tm := util.Marshal(tmsg)
	if execType == 0 {
		tinfo, err := GetWsInfo(msg.TargetClientId)
		if err != nil {
			fmt.Printf("err:%+v\n", err)
			return
		}
		tinfo.send(tm)
	} else {
		logsClis.clients.Range(func(key, value interface{}) bool {
			k := key.(types.ClientId)
			v := value.(*model.ExecutorInfo)
			if (v.Pid == 0 && v.ClientId == info.ClientId) || (v.Pid != 0 && v.Pid == msg.ExecMessage.Pid && v.ClientId == info.ClientId) {
				tinfo, err := GetWsInfo(k)
				if err != nil {
					fmt.Printf("err:%+v\n", err)
				} else {
					tinfo.send(tm)
				}
			}
			return true
		})
	}
}

func handleClientBeginExec(info *WsInfo, msg *model.Message) {
	execClis.Add(info.ClientId, &model.ExecutorInfo{
		ClientId: info.ClientId,
		Pid:      msg.ExecMessage.Pid,
		Command:  msg.ExecMessage.Command,
	})
	tmsg := &model.Message{
		MessageType:    constants.COMMON,
		Message:        msg.Message,
		TargetClientId: info.ClientId,
		ExecMessage: &model.ExecutorMessage{
			Ip:      info.Ip,
			Command: msg.ExecMessage.Command,
			Pid:     msg.ExecMessage.Pid,
		},
	}
	tm := util.Marshal(tmsg)
	tinfo, err := GetWsInfo(msg.TargetClientId)
	if err != nil {
		fmt.Printf("err:%+v\n", err)
		return
	}
	tinfo.send(tm)
}

func handleClientEndExec(info *WsInfo, msg *model.Message) {
	execClis.DeleteInfo(info.ClientId, msg.ExecMessage.Pid)
	deleteZdoraKeys := logsClis.DeleteInfo(info.ClientId, msg.ExecMessage.Pid)
	for _, key := range deleteZdoraKeys {
		tinfo, err := GetWsInfo(key)
		if err != nil {
			fmt.Printf("err:%+v\n", err)
		} else {
			tinfo.send(util.Marshal(&model.Message{
				MessageType: constants.CLOSE,
			}))
		}
	}
}

func handleLogs(info *WsInfo, msg *model.Message) {
	if msg.ExecMessage.Pid == 0 {
		if _, ok := execClis.clients.Load(msg.TargetClientId); ok {
			logsClis.Add(info.ClientId, &model.ExecutorInfo{
				ClientId: msg.TargetClientId,
				Pid:      0,
			})
		} else {
			info.send(util.Marshal(&model.Message{
				MessageType: constants.CLOSE,
				Message:     "客户端" + string(msg.TargetClientId) + "处于空闲状态",
			}))
		}
	} else {
		einfo := execClis.Get(msg.TargetClientId, msg.ExecMessage.Pid)
		if einfo == nil {
			info.send(util.Marshal(&model.Message{
				MessageType: constants.CLOSE,
				Message:     "客户端" + string(msg.TargetClientId) + "进程" + string(msg.ExecMessage.Pid) + "没有运行",
			}))
			return
		} else {
			logsClis.Add(info.ClientId, &model.ExecutorInfo{
				ClientId: msg.TargetClientId,
				Pid:      msg.ExecMessage.Pid,
			})
		}
	}
}
