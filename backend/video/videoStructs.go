package video

import "github.com/gorilla/websocket"

var VideoUpgrader = websocket.Upgrader{

	ReadBufferSize: 1024,

	WriteBufferSize: 1024,
}
