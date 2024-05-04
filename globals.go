package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

var (
	connections map[string]*websocket.Conn
	lock        sync.RWMutex
)

func init() {
	connections = make(map[string]*websocket.Conn)
}
