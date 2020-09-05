package err

import (
	"fmt"
	"zdora/types"
)

type ErrorNotFound struct {
	ClientId types.ClientId
}

func (e ErrorNotFound) Error() string {
	return fmt.Sprintf("连接不存在:[%s]", e.ClientId)
}
