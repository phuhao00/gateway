package main

import (
	"github.com/phuhao00/fuse"
	"github.com/phuhao00/network"
)

type ClientServer struct {
	real *network.Server
	//与其他服绑定信息
	FromInnerCh chan interface{}
	ToInnerCh   chan interface{}
	router      *fuse.Router
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
	//todo check
	s.ToInnerCh <- packet
}

func (s *ClientServer) Router(data interface{}) {
	handler := s.router.GetHandler(data.(*network.Packet))
	if handler != nil {
		val := data.(*network.Packet)
		handler(val, nil) //todo
	}
}

func (s *ClientServer) Register() {
	s.router.AddRoute(333, s.RegisterLoginInfo)
	s.router.AddRoute(444, s.ForwardServerPacket)
}
