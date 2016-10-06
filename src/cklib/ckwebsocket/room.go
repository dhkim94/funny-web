package ckwebsocket

import (
	"cklib/env"
)

type Room struct {
	// room 의 이름
	name string

	// 동일한 room 에 있는 client 목록
	clients map[*Client]bool

	// broadcast 할 메시지
	BroadcastMsg chan []byte

	// room 에 join 하는 client
	JoinClient chan *Client

	// room 을 떠나는 client
	LeaveClient chan *Client
}

// 새로운 room 을 생성 한다.
// name : room name
func NewRoom(name string) *Room {
	return &Room{
		name: name,
		clients: make(map[*Client]bool),
		BroadcastMsg: make(chan []byte),
		JoinClient: make(chan *Client),
		LeaveClient: make(chan *Client),
	}
}

func (room *Room) Run() {
	slog := env.GetLogger()

	for {
		select {
		case client := <-room.JoinClient:
			room.clients[client] = true
			slog.Info("join client:%d in room:%s", client.Id, room.name)
		case client := <-room.LeaveClient:
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				close(client.SendMsg)
				slog.Info("leave client:%d from room:%s", client.Id, room.name)
			}
		case msg := <-room.BroadcastMsg:
			for client := range room.clients {
				select {
				case client.SendMsg <-msg:
				default:
					close(client.SendMsg)
					delete(room.clients, client)
					slog.Info("leave client:%d from room:%s", client.Id, room.name)
				}
			}
		}
	}
}