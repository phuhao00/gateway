package main

import (
	"github.com/phuhao00/network"
	"sync"
)

type InnerServer struct {
	real     *network.Server
	sessions sync.Map //*network.Session
	//哪些消息交给这个服处理（messageId 范围控制）
	FromClientCh chan interface{}
	ToClientCh   chan interface{}
}

func NewInnerServer() *InnerServer {
	return &InnerServer{
		real: network.NewServer(""),
	}
}

func (s *InnerServer) loop() {
	s.real.OnSessionPacket = s.MessageHandler
	s.real.Run()
	for {
		select {
		case data := <-s.FromClientCh:
			s.Router(data)
		}
	}
}

func (s *InnerServer) AddSession(session *network.Session) {
	s.sessions.Store(session, struct{}{})
}

func (s *InnerServer) DeleteSession(session *network.Session) {
	s.sessions.Delete(session)

}

func (s *InnerServer) MessageHandler(packet *network.SessionPacket) {

}

func (s *InnerServer) Router(interface{}) {
	//todo 发送给对应的服务器处理
}
