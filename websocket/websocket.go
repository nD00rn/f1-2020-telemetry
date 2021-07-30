package websocket

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan *string)
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
    }

    // register client
    clients[ws] = true
}

func BroadcastMessage(message string) {
    broadcast <- &message
}

func Broadcast() {
    for {
        val := <-broadcast

        // send to every client that is currently connected
        for client := range clients {
            err := client.WriteMessage(
                websocket.TextMessage,
                []byte(*val),
            )
            if err != nil {
                log.Printf("Websocket error: %s", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}
