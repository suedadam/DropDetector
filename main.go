package main

import (
	"os"
	"log"
	"encoding/json"
	"time"
	"github.com/deckarep/gosx-notifier"

	"github.com/geNAZt/minecraft-status/data"
	"github.com/geNAZt/minecraft-status/protocol"
)

var oldplayer int

func main() {
    defer func() {
        r := recover()
        if r != nil {
        	note := gosxnotifier.NewNotification("Program crashed, recovering!")
			note.Title = "Drop detector"
			note.Push()
        }
    }()
	addr := os.Args[1]
	for {
		time.Sleep(2 * time.Second);
		query(addr)
		log.Printf("Players: %d\n", oldplayer)
	}
}

func query(ip string) {

	conn, err := protocol.NewNetClient(ip)
	if err != nil {
		note := gosxnotifier.NewNotification("Failed to connect!")
		note.Title = "Drop detector"
		note.Push()
	}
	defer conn.Close()

	conn.SendHandshake()
	conn.State = protocol.Status

	conn.SendStatusRequest()
	statusPacket, err := conn.ReadPacket()
	if err != nil {
		log.Printf("%s\n", err)
	}
	status := &data.Status{}
	_ = json.Unmarshal([]byte(statusPacket.(protocol.StatusResponse).Data), status)

	newplayer := int(status.Players.Online)
	if oldplayer != 0 {
		if newplayer <= oldplayer - 20 {
			note := gosxnotifier.NewNotification("DROP!!")
			note.Title = "Drop detector"
			note.Push()
		}
	}
	oldplayer = newplayer

}
