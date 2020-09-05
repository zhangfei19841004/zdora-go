package util

import (
	"sync"
	"zdora/types"
)

type ClientsMap struct {
	clients sync.Map
}

func (cm *ClientsMap) GetClients() sync.Map {
	return cm.clients
}

func (cm *ClientsMap) Add(clientId types.ClientId, info interface{}) {
	cm.clients.Store(clientId, info)
}

func (cm *ClientsMap) Delete(clientId types.ClientId) {
	cm.clients.Delete(clientId)
}
