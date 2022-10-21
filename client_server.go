package main

import (
	"github.com/phuhao00/network"
)

type ClientServer struct {
	real *network.Server
	//与其他服绑定信息
	FromInnerCh chan interface{}
	ToInnerCh   chan interface{}
}

func NewClientServer() *ClientServer {
	return &ClientServer{
		//real:     network.NewServer(""),
	}
}

func (s *ClientServer) loop() {
	s.real.MessageHandler = s.MessageHandler
	for {
		select {
		case data := <-s.FromInnerCh:
			s.Router(data)
		}
	}

}

func (s *ClientServer) MessageHandler(packet *network.Packet) {

}

func (s *ClientServer) Router(interface{}) {
	//todo  发送给对应客户端
}
