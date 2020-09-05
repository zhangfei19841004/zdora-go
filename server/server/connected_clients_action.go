package server

import (
	"zdora/err"
	"zdora/model"
	"zdora/types"
)

func GetWsInfo(clientId types.ClientId) (*WsInfo, error) {
	clients := clients.GetClients()
	if info, ok := clients.Load(clientId); ok {
		return info.(*WsInfo), nil
	}
	return nil, err.ErrorNotFound{ClientId: clientId}
}

func ListAllClients() []*model.PsCommandInfo {
	var infos []*model.PsCommandInfo
	clients := clients.GetClients()
	clients.Range(func(key, value interface{}) bool {
		info := value.(*WsInfo)
		if !info.IsCobraClient() {
			if v, ok := execClis.clients.Load(key); ok {
				einfos := v.([]*model.ExecutorInfo)
				for _, einfo := range einfos {
					infos = append(infos, &model.PsCommandInfo{
						ClientId: info.ClientId,
						Ip:       info.Ip,
						Pid:      einfo.Pid,
						Commond:  einfo.Command,
					})
				}
			} else {
				infos = append(infos, &model.PsCommandInfo{
					ClientId: info.ClientId,
					Ip:       info.Ip,
					Pid:      0,
					Commond:  []string{},
				})
			}
		}
		return true
	})
	return infos
}
