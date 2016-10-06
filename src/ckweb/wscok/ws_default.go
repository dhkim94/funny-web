package wsock

import (
	"net/http"
	"github.com/gorilla/websocket"
	"cklib/env"
	"cklib/ckwebsocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func WsDefaultConnect(rootRoom *ckwebsocket.Room, clientSn uint64, w http.ResponseWriter, r *http.Request) {
	slog := env.GetLogger()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Err("failed client connected to ws default")
		return
	}

	client := &ckwebsocket.Client{
		Id: clientSn,
		RootRoom: rootRoom,
		Conn: conn,
		SendMsg: make(chan []byte, 256),
	}

	slog.Info("connected client:%d in ws default", client.Id)

	// client 를 웹소켓 주소에 만들어 지는 root room 에 입장 시킨다.
	// 추후 접근한 모두에게 broadcast 하고 싶다면 root room 에 들어 있는 모든 사람에게 메시지를 발송 하면 된다.
	rootRoom.JoinClient <- client

	// client 에 메시지를 보내는 루틴을 실행
	go client.Write()

	// client 에서 메시지가 들어오면 읽는 루틴을 실행
	client.Read()


}
