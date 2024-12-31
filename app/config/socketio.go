package config

import (
	"github.com/gofiber/contrib/socketio"
)

var SocketIOClient = make(map[string]map[string]*socketio.Websocket)

// Menambahkan koneksi WebSocket ke dalam map berdasarkan userId
// Menambahkan koneksi WebSocket ke dalam map berdasarkan route dan userId
func AddSocketIOClient(route, userId string, kws *socketio.Websocket) {
	if _, exists := SocketIOClient[route]; !exists {
		SocketIOClient[route] = make(map[string]*socketio.Websocket)
	}
	SocketIOClient[route][userId] = kws
}

// Mengirim notifikasi ke semua client di route tertentu
func NotifyAllSocketIOClient(route, message string) {
	if clients, exists := SocketIOClient[route]; exists {
		for _, client := range clients {
			client.Emit([]byte(message), socketio.TextMessage)
		}
	}
}

// Mengirimkan pesan ke client berdasarkan userId
func NotifySocketIOClientIn(route, userId, message string) {
	if clients, exists := SocketIOClient[route]; exists {
		if client, ok := clients[userId]; ok {
			client.Emit([]byte(message), socketio.TextMessage)
		}
	}
}
