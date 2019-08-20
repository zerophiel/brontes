package websocket

import (
    "sync"

    "github.com/graarh/golang-socketio"
    "github.com/graarh/golang-socketio/transport"
)

var (
    once     sync.Once
)

func WebSocketInit () {
    once.Do(func () {
        //connect to server, you can use your own transport settings
        _, err := gosocketio.Dial(
            gosocketio.GetUrl("localhost", 80, false),
            transport.GetDefaultWebsocketTransport(),
        )
        if err != nil {

        }
        //do something, handlers and functions are same as server ones

        //close connection
        // c.Close()
    })
}
