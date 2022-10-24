package main

import (
	"errors"
	"github.com/phuhao00/fuse"
	"github.com/phuhao00/network"
	"math/rand"
	"sort"
	"sync"
)

type InnerServer struct {
	real         *network.Server
	servers      sync.Map
	FromClientCh chan interface{}
	ToClientCh   chan interface{}
	router       *fuse.Router
}

type ServerInfo struct {
	ServerAddr   string
	ServerID     string
	ServerType   uint32
	ProIndex     uint32
	ZoneId       int
	HandleMsgIds []uint16
	conn         *network.TcpConnX
	players      int64
	maxPlayer    int64
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

func (s *InnerServer) getServerInfoList() []*ServerInfo {
	sis := make([]*ServerInfo, 0, 128)
	s.servers.Range(func(key, value interface{}) bool {
		si, ok := value.(*ServerInfo)
		if ok {
			sis = append(sis, si)
		}
		return true
	})

	if len(sis) > 0 {
		sort.Slice(sis[:], func(i, j int) bool {
			return sis[i].ProIndex < sis[j].ProIndex
		})
	}

	return sis
}

func (s *InnerServer) addServerInfo(srvID, address string, msgIds []uint32, conn *network.TcpConnX, proIndex uint32, zoneId int) {
	old, ok := s.servers.Load(srvID)
	if ok {
		oldS, _ := old.(*ServerInfo)
		if oldS.conn != conn {
			oldS.conn.Close()
		}
	}
	si := &ServerInfo{}

	for _, msgID := range msgIds {
		si.HandleMsgIds = append(si.HandleMsgIds, uint16(msgID))
	}
	si.ServerAddr = address
	si.ServerID = srvID
	si.ServerType = 111 //todo
	si.ProIndex = proIndex
	si.ZoneId = zoneId
	s.servers.Store(srvID, si)
}

func (s *InnerServer) removeServerInfo(srvId string) {
	si, ok := s.servers.Load(srvId)
	if !ok {
		return
	}
	s.servers.Delete(srvId)
	serverInfo := si.(*ServerInfo)
	GetClientServerInstance().onlineServerDisconnected(srvId, serverInfo.ServerAddr, serverInfo.ZoneId, serverInfo.ProIndex)
}

func (s *InnerServer) getServerInfo(srvID string) (*ServerInfo, error) {
	si, ok := s.servers.Load(srvID)
	if !ok {
		return nil, errors.New("do not exist")
	}
	return si.(*ServerInfo), nil
}

func (s *InnerServer) sendMsgToServer(srvID string, cmd uint16, msg interface{}) bool {
	si, ok := s.servers.Load(srvID)
	if !ok {
		return false
	}
	return si.(*ServerInfo).conn.AsyncSendLastPacket(cmd, msg)
}

func (s *InnerServer) HandlerMsg(srvID string, userID uint64, msgID uint16, data []byte) {
	si, ok := s.servers.Load(srvID)
	if !ok {
		return
	}
	si.(*ServerInfo).conn.AsyncSendLastPacket(msgID, data)
}

func (s *InnerServer) getOptimalServer() (string, bool) {
	idle := make([]*ServerInfo, 0)
	busy := make([]*ServerInfo, 0)
	hot := make([]*ServerInfo, 0)
	for _, si := range s.getServerInfoList() {
		rate := (float32(si.players) / float32(si.maxPlayer)) * 100
		if rate < 100 && rate > 80 {
			hot = append(hot, si)
		} else if rate <= 80 && rate > 40 {
			busy = append(busy, si)
		} else if rate <= 40 {
			idle = append(idle, si)
		}
	}
	if len(busy) > 0 {
		return busy[rand.Int()%len(busy)].ServerID, true
	}
	if len(idle) > 0 {
		return idle[rand.Int()%len(idle)].ServerID, true
	}
	if len(hot) > 0 {
		return hot[rand.Int()%len(hot)].ServerID, true
	}
	return "", false
}
