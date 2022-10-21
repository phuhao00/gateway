package main

import (
	"github.com/phuhao00/network"
	"sync"
)

type InnerServer struct {
	real        *network.Server
	serversInfo sync.Map //*network.Session
	//哪些消息交给这个服处理（messageId 范围控制）
	FromClientCh chan interface{}
	ToClientCh   chan interface{}
}

func NewInnerServer() *InnerServer {
	return &InnerServer{
		//real: network.NewServer("127.0.0.1:3456"),
	}
}

func (s *InnerServer) loop() {
	s.real.MessageHandler = s.MessageHandler
	s.real.Run()
	for {
		select {
		case data := <-s.FromClientCh:
			s.Router(data)
		}
	}
}

func (s *InnerServer) MessageHandler(packet *network.Packet) {

}

func (s *InnerServer) Router(interface{}) {

	//todo 发送给对应的服务器处理
}
