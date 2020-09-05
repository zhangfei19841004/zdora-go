package server

import (
	"sync"
	"zdora/model"
	"zdora/server/util"
	"zdora/types"
)

var clients = &util.ClientsMap{} // connected clients

var execClis = &execClients{}

var logsClis = &logsClients{} //key是zdora的clientId，[]里的ExecutorInfo里的key是exec里的clientId

type execClients struct {
	clients sync.Map
}

func (exec *execClients) Add(clientId types.ClientId, info *model.ExecutorInfo) {
	var infos []*model.ExecutorInfo
	if v, ok := exec.clients.Load(clientId); ok {
		infos = v.([]*model.ExecutorInfo)
	}
	infos = append(infos, info)
	exec.clients.Store(clientId, infos)
}

func (exec *execClients) Get(clientId types.ClientId, pid types.PID) *model.ExecutorInfo {
	if v, ok := exec.clients.Load(clientId); ok {
		infos := v.([]*model.ExecutorInfo)
		for _, info := range infos {
			if info.ClientId == clientId && info.Pid == pid {
				return info
			}
		}
	}
	return nil
}

func (exec *execClients) DeleteInfo(clientId types.ClientId, pid types.PID) {
	if pid == 0 {
		exec.clients.Delete(clientId)
		return
	}
	var infos []*model.ExecutorInfo
	if value, ok := exec.clients.Load(clientId); ok {
		temps := value.([]*model.ExecutorInfo)
		for _, temp := range temps {
			if temp.ClientId != clientId || temp.Pid != pid {
				infos = append(infos, temp)
			}
		}
	}
	if len(infos) == 0 {
		exec.clients.Delete(clientId)
	} else {
		exec.clients.Store(clientId, infos)
	}
}

//zdoraClientId为key, value中如果pid为0，则logs该client下的所有pid log,如果pid不为0，则Logs该client下对应的pid的log
type logsClients struct {
	clients sync.Map
}

func (logs *logsClients) Add(zdoraClientId types.ClientId, info *model.ExecutorInfo) {
	logs.clients.Store(zdoraClientId, info)
}

func (logs *logsClients) DeleteInfo(clientId types.ClientId, pid types.PID) []types.ClientId {
	var deleteZdoraKeys []types.ClientId
	logs.clients.Range(func(key, value interface{}) bool {
		v := value.(*model.ExecutorInfo)
		if pid == 0 {
			if v.ClientId == clientId {
				deleteZdoraKeys = append(deleteZdoraKeys, key.(types.ClientId))
				logs.clients.Delete(key)
			}
		} else {
			if v.Pid != 0 && pid == v.Pid && v.ClientId == clientId {
				deleteZdoraKeys = append(deleteZdoraKeys, key.(types.ClientId))
				logs.clients.Delete(key)
			} else if v.Pid == 0 && v.ClientId == clientId {
				if _, ok := execClis.clients.Load(clientId); !ok {
					deleteZdoraKeys = append(deleteZdoraKeys, key.(types.ClientId))
					logs.clients.Delete(key)
				}
			}
		}
		return true
	})
	return deleteZdoraKeys
}

func (logs *logsClients) Delete(zdoraClientId types.ClientId) {
	logs.clients.Delete(zdoraClientId)
}
