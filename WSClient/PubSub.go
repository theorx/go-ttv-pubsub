package WSClient

import (
	"errors"
	"ttvWS/Topic"
)

var ErrorOperationFailed = errors.New("sub/unsub operation failed")
var ErrorNotConnected = errors.New("not connected")

func (c *Client) Subscribe(topics []Topic.Topic) error {
	if c.isConnected() == false {
		return ErrorNotConnected
	}

	resultFN, err := c.Request(&OutgoingMessage{
		Type: "LISTEN",
		Data: struct {
			Topics    []Topic.Topic `json:"topics,omitempty"`
			AuthToken string        `json:"auth_token,omitempty"`
		}{
			Topics:    topics,
			AuthToken: c.authToken,
		},
	})

	if err != nil {
		return err
	}

	result := resultFN()

	if len(result.Error) == 0 {
		for _, t := range topics {
			exists := false
			for _, ct := range c.topics {
				if t == ct {
					exists = true
				}
			}

			if exists == false {
				c.topics = append(c.topics, t)
			}
		}

		return nil
	}

	return ErrorOperationFailed
}

func (c *Client) Unsubscribe(topics []Topic.Topic) error {
	if c.isConnected() == false {
		return ErrorNotConnected
	}

	resultFN, err := c.Request(&OutgoingMessage{
		Type: "UNLISTEN",
		Data: struct {
			Topics    []Topic.Topic `json:"topics,omitempty"`
			AuthToken string        `json:"auth_token,omitempty"`
		}{
			Topics:    topics,
			AuthToken: c.authToken,
		},
	})

	if err != nil {
		return err
	}

	result := resultFN()

	if len(result.Error) == 0 {
		//remove topics from list
		newList := make([]Topic.Topic, 0)

		for _, t := range c.topics {
			exists := false
			for _, p := range topics {
				if t == p {
					exists = true
				}
			}

			if exists == false {
				newList = append(newList, t)
			}
		}

		c.topics = newList
		return nil
	}

	return ErrorOperationFailed
}
