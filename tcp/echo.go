package tcp

import (
	"bufio"
	"context"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

/**
 * @Description
 * @Author ZzzWw
 * @Date 2022-06-27 9:49
 **/

// 客户端的实体

type EchoClient struct {
   Conn net.Conn
   Waiting wait.Wait
}

func (e *EchoClient) Close() error{
	e.Waiting.WaitWithTimeout(10*time.Second) //等待10秒
	_ = e.Conn.Close()  //关闭
	return nil
}
type EchoHandler struct {
	activeConn sync.Map
	closing atomic.Boolean //判断是否在关闭
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {

	if handler.closing.Get(){
		_ = conn.Close()
	}
	client:=&EchoClient{
		Conn: conn,
	}
	handler.activeConn.Store(client, struct {}{})
	reader:=bufio.NewReader(conn)
    for{
		msg,err:=reader.ReadString('\n')
		if err!=nil{
			if err==io.EOF{
				logger.Info("Connecting close")
				handler.activeConn.Delete(client)
			}else{
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		b:=[]byte(msg)
		_,_ = conn.Write(b)
		client.Waiting.Done()
	}
}

func (handler *EchoHandler) Close() error {
        logger.Info("handler shutting down")
		handler.closing.Set(true)
		handler.activeConn.Range(func(key, value interface{}) bool {
			client:=key.(*EchoClient)
			_ = client.Conn.Close()
			return true
		})

       return nil
}


