package tcp

import (
	"context"
	"net"
)

/**
 * @Description
 * @Author ZzzWw
 * @Date 2022-06-24 10:28
 **/

// Handler 业务逻辑处理接口
type Handler interface {
	Handle(ctx context.Context,conn net.Conn)
	Close() error
}