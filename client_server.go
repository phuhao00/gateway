package main

import (
	"github.com/phuhao00/network"
	"sync"
)

type ClientServer struct {
	real     *network.Server
	sessions sync.Map //*network.Session
	//与其他服绑定信息
	FromInnerCh chan interface{}
	ToInnerCh   chan interface{}
}

func NewClientServer() *ClientServer {
	return &ClientServer{
		real:     network.NewServer(""),
		sessions: sync.Map{},
	}
}

func (s *ClientServer) loop() {
	s.real.OnSessionPacket = s.MessageHandler
	for {
		select {
		case data := <-s.FromInnerCh:
			s.Router(data)
		}
	}

}

func (s *ClientServer) AddSession(session *network.Session) {
	s.sessions.Store(session, struct{}{})
}

func (s *ClientServer) DeleteSession(session *network.Session) {
	s.sessions.Delete(session)

}

func (s *ClientServer) MessageHandler(packet *network.SessionPacket) {

}

func (s *ClientServer) Router(interface{}) {
	//todo  发送给对应客户端
}
