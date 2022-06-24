package tcp

import (
	"context"
	"go-redis/interface/tcp"
	"go-redis/lib/logger"
	"net"
	"sync"
)

/**
 * @Description
 * @Author ZzzWw
 * @Date 2022-06-24 10:31
 **/


type Config struct {
	Address string //监听地址

}

func ListenAndServeWithSignal(cfg *Config,handler tcp.Handler)error{
	listener, err := net.Listen("tcp", cfg.Address)
	closeChan:=make(chan struct{})
	if err!=nil{
		return err
	}
	logger.Info("start listen")
	ListenAndServer(listener,handler,closeChan)
	return nil
}

func ListenAndServer(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}){
	ctx := context.Background()
	var waitDone sync.WaitGroup //等待所有客户端退出
	//不断接受新链接
	for true{
		conn, err := listener.Accept()
		//出错，则跳出循环
		if err!=nil{
			break
		}
		logger.Info("accepted link") //有新的链接
		waitDone.Add(1) //waitGroup + 1
		//新建协程 处理业务
		go func() {
			defer func() {
				waitDone.Done() //waitGroup -1
			}()
			handler.Handle(ctx,conn)
		}()
	}
	waitDone.Wait() //等待退出
}