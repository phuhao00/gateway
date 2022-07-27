package main

import "github.com/phuhao00/network"

type ServerClient struct {
	session  *network.Session
	category int //服务器类型
	//哪些消息交给这个服处理（messageId 范围控制）
	playerCounter int32 //玩家数量
}
