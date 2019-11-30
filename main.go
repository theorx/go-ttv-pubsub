package main

import (
	"log"
	"time"
	"ttvWS/Topic"
	"ttvWS/WSClient"
)

const TheOrXPlays = 64417816

func main() {
	log.Println("Running ttvWS")

	log.SetFlags(log.Lshortfile | log.Ldate)

	c, _ := WSClient.CreateClient("")

	_ = c.Subscribe([]Topic.Topic{
		Topic.ModerationAction(TheOrXPlays, TheOrXPlays),
		Topic.Bits(TheOrXPlays),
	})

	c.SetModerationHandler(func(message WSClient.ModerationActionMsg) {

		log.Println("Moderation action much?", message)

	})

	c.SetBitsHandler(func(message WSClient.BitsMsg) {
		log.Println("Bit handler", message)
	})

	time.Sleep(time.Second * 1000)

}
