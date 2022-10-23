package main

import (
	"github.com/phuhao00/fuse"
	"github.com/phuhao00/network"
	"sync"
)

type InnerServer struct {
	real        *network.Server
	serversInfo sync.Map //*network.Session
	//哪些消息交给这个服处理（messageId 范围控制）
	FromClientCh chan interface{}
	ToClientCh   chan interface{}
	router       *fuse.Router
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
	//如果是注册节点信息

	s.ToClientCh <- packet
}

func (s *InnerServer) Router(data interface{}) {
	//world server 多节点支持
	//一般逻辑都会经过world server 做转发
	//战斗服的话，会直连客户端

	//get serverUId
	handler := s.router.GetHandler(data.(*network.Packet))
	if handler != nil {
		val := data.(*network.Packet)
		handler(val, nil)
	}
	//todo 发送给对应的服务器处理
}

func (s *InnerServer) Register() {
	s.router.AddRoute(11, s.ServerInfoRegister)
	s.router.AddRoute(22, s.ForwardClientPacket)
}
