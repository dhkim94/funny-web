package ckwebsocket

import (
	"github.com/gorilla/websocket"
	"time"
	"cklib/env"
	"bytes"
)

const (
	// client 에 허가된 최대 메시지 사이즈
	maxMessageSize	= 512

	// client 에게 데이터 write 하는 대기 시간
	// 대기 시간 동안 write 하지 못하면 client 와의 접속을 종료 시킨다.
	writeWait	= 10 *time.Second

	// client 에서 받는 pong 의 다음 pong 까지의 최대 시간 간격
	// client 에서 설정한 pong 시간 가격내에 어떠한 메시지가 오지 않으면 접속 종료 시킨다.
	pongWait	= time.Second * 60

	// client 에 ping 을 보내는 시간 간격
	// pongWait 보다 작은 값이어야 한다.
	pingPeriod	= (pongWait * 9) / 10
)

var (
	newline		= []byte{'\n'}
	space		= []byte{' '}
)

type Client struct {
	// 접근한 web socket url 의 root room
	RootRoom *Room

	// 접속된 client 의 connection
	Conn *websocket.Conn

	// 서버가 client 에게 보낼 메시지
	SendMsg chan []byte

	// client serial number
	Id uint64
}

// client 로 부터 받은 메시지를 읽는다.
func (client *Client) Read() {
	slog := env.GetLogger()
	slog.Info("start read routine for client:%d", client.Id)

	defer func() {
		client.RootRoom.LeaveClient <-client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil; })

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				slog.Warn("read message error from client:%d. error: %v", err, client.Id)
			}
			break
		}

		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		slog.Info("> client:%d recv msg[%s]", client.Id, msg)

		// todo message 받은 것에 맞는 행동을 해야 한다.

	}
}

func (client *Client) Write() {
	slog := env.GetLogger()
	slog.Info("start write routine for client:%d", client.Id)

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.SendMsg:
			if !ok {
				// client 로 부터 메시지 받는 것에 문제가 있다면 클라이언트에게 접속 종료 명령을 보낸다.
				slog.Info("send close command to client:%d", client.Id)
				client.sendToClient(websocket.CloseMessage, []byte{})
				return
			}

			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Warn("failed get data sender-writer to client:%d. ignore data send to client",
					client.Id)
				return
			}
			w.Write(msg)

			// todo 여기서 무슨 add queued 가 필요하지 ???

			if err := w.Close(); err != nil {
				slog.Warn("failed close data sender-writer to client:%d. ignore data send to client",
					client.Id)
				return
			}
		case <-ticker.C:
			if err := client.sendToClient(websocket.PingMessage, []byte{}); err != nil {
				slog.Warn("failed send ping command to client:%d", client.Id)
				return
			}
			slog.Debug("send ping command to client:%d", client.Id)
		}
	}
}

func (client *Client) sendToClient(msgType int, payload []byte) error {
	client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return client.Conn.WriteMessage(msgType, payload)
}