package TTVClient

import (
	"encoding/json"
	"errors"
	"github.com/theorx/go-ttv-pubsub/pkg/Topic"
	"time"
)

func (c *Client) readLoop() {

	for {
		message := &IncomingMessage{}
		err := c.conn.ReadJSON(&message)

		if err != nil {
			c.log("Closed read loop due to read failure", err)

			c.reconnect()
			return
		}
		c.handleCatchAll(*message)

		if ok, err := c.handleControl(*message); ok {
			//if error, then we exit the function to shut down
			if err != nil {
				return
			}
			continue
		}

		if c.handleNonce(*message) {
			continue
		}

		if c.handleTopics(*message) {
			continue
		}

		c.handleUnknown(*message)
	}

}

func (c *Client) handleNonce(msg IncomingMessage) bool {

	if val, ok := c.inFlight[msg.Nonce]; ok {
		val <- msg
		return true
	}

	return false
}

func (c *Client) handleTopics(msg IncomingMessage) bool {

	if len(msg.Data.Topic) == 0 {
		return false
	}

	switch Topic.GetType(msg.Data.Topic) {
	case Topic.TypeBits:
		m := &BitsMsg{}

		err := json.Unmarshal([]byte(msg.Data.Message), &m)
		if err != nil {
			return false
		}

		if c.bitsHandler != nil {
			c.bitsHandler(*m)
		}
		return true
	case Topic.TypeBitsBadgeNotification:
		m := &BitsBadgeMsg{}

		err := json.Unmarshal([]byte(msg.Data.Message), &m)
		if err != nil {
			return false
		}

		if c.bitsBadgeHandler != nil {
			c.bitsBadgeHandler(*m)
		}
		return true
	case Topic.TypeSubscriptions:
		m := &SubscriptionMsg{}

		err := json.Unmarshal([]byte(msg.Data.Message), &m)
		if err != nil {
			return false
		}

		if c.subscriptionsHandler != nil {
			c.subscriptionsHandler(*m)
		}
		return true
	case Topic.TypeCommerce:
		m := &CommerceMsg{}

		err := json.Unmarshal([]byte(msg.Data.Message), &m)
		if err != nil {
			return false
		}

		if c.commerceHandler != nil {
			c.commerceHandler(*m)
		}
		return true
	case Topic.TypeWhispers:
		m := &WhisperMsg{}

		err := json.Unmarshal([]byte(msg.Data.Message), &m)
		if err != nil {
			return false
		}

		if c.whispersHandler != nil {
			c.whispersHandler(*m)
		}
		return true
	case Topic.TypeModerationAction:
		m := &ModerationActionMsg{}

		err := json.Unmarshal([]byte(msg.Data.Message), &m)
		if err != nil {
			return false
		}

		if c.moderationHandler != nil {
			c.moderationHandler(*m)
		}
		return true
	default:
		return false
	}

}

func (c *Client) handleControl(msg IncomingMessage) (bool, error) {
	msgType := msg.Type

	if Topic.Type(msgType) == Topic.TypePong {
		c.lastPongTS = time.Now().Unix()
		return true, nil
	}

	if Topic.Type(msgType) == Topic.TypeReconnect {
		c.log("Reconnect msg received..")
		c.reconnect()
		return true, errors.New("shutdown")
	}

	return false, nil
}

func (c *Client) handleUnknown(msg IncomingMessage) {
	if c.unknownHandler != nil {
		c.unknownHandler(msg)
		return
	}
	c.log("Unknown message received, set client.SetUnknownHandler(handler) to handle it")
}

func (c *Client) handleCatchAll(msg IncomingMessage) {
	if c.catchAllHandler != nil {
		c.catchAllHandler(msg)
	}
}
