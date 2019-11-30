package TTVClient

import (
	"sync/atomic"
	"time"
)

func (c *Client) pingLoop() {
	for {

		if atomic.LoadInt64(&c.connectionStatus) == 0 {
			c.log("Connection status 0, shutting down pingLoop")
			return
		}

		if time.Now().Unix()%60 == 0 {
			c.ping()
		}

		if time.Now().Unix()%10 == 0 {

			if c.lastPingTS < (time.Now().Unix() - 120) {
				//shit hits the fan and we're most likely not connected to the service anymore, handle reconnection?!? todo:@@@
				c.log("Last ping was too long ago, 120+ seconds, reconnecting!")
				c.reconnect()
				return
			}

			if c.lastPongTS-c.lastPingTS > 20 {
				c.log("Last pong was more than 20 seconds late, reconnecting!")
				c.reconnect()
				return
			}
		}

		time.Sleep(time.Second * 1)
	}
}

func (c *Client) ping() {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	err := c.conn.WriteJSON(&OutgoingMessage{
		Type: "PING",
	})

	if err == nil {
		c.lastPingTS = time.Now().Unix()
	}
}
