package WSClient

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"sync/atomic"
	"time"
	"ttvWS/Topic"
)

const TwitchSecureWSHost = "wss://pubsub-edge.twitch.tv/"

type Client struct {
	endpoint  string
	conn      *websocket.Conn
	authToken string
	/* Map on requests in flight with response channels */
	inFlight         map[string]chan IncomingMessage
	connectionStatus int64
	topics           []Topic.Topic
	//todo implement log function in the future

	err error

	writeLock  *sync.Mutex
	lastPingTS int64
	lastPongTS int64

	//topic handlers
	bitsHandler          BitsHandlerFunction
	bitsBadgeHandler     BitsBadgeHandlerFunction
	subscriptionsHandler SubscriptionsHandlerFunction
	commerceHandler      CommerceHandlerFunction
	whispersHandler      WhispersHandlerFunction
	moderationHandler    ModerationActionHandlerFunction
	//other handlers
	catchAllHandler HandlerFunction
	unknownHandler  HandlerFunction
}

func CreateClient(authToken string) (*Client, error) {

	client := &Client{
		endpoint:   TwitchSecureWSHost,
		authToken:  authToken,
		writeLock:  &sync.Mutex{},
		lastPingTS: time.Now().Unix(),
		lastPongTS: time.Now().Unix(),
		inFlight:   make(map[string]chan IncomingMessage),
	}

	return client, client.init()
}

func (c *Client) init() error {
	//set the connection status
	atomic.StoreInt64(&c.connectionStatus, 1)
	c.reconnect()

	return c.err
}

func (c *Client) isConnected() bool {
	return atomic.LoadInt64(&c.connectionStatus) == 1
}

func (c *Client) reconnect() {

	if atomic.CompareAndSwapInt64(&c.connectionStatus, 1, 0) == false {
		return
	}

	c.err = nil

	//close all the channels in flight
	for _, ch := range c.inFlight {
		close(ch)
	}
	//reset the map
	c.inFlight = make(map[string]chan IncomingMessage)

	if c.conn != nil {
		err := c.conn.Close()

		if err != nil {
			log.Println("Failed closing old connection", err)
		}
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.endpoint, nil)
	c.conn = conn

	if err != nil {
		c.err = err
		return
	}
	atomic.StoreInt64(&c.connectionStatus, 1) //connection back up

	go c.readLoop()
	go c.pingLoop()
	c.ping()

	//call listen with the stored topics -- first back up the values and then save them after the Subscribe() call to prevent having duplicate values
	if len(c.topics) > 0 {
		c.err = c.Subscribe(c.topics)
	}
}

func (c *Client) generateNonce() (string, chan IncomingMessage) {
	id, _ := uuid.NewUUID() //ignore error, this error can never happen, error is always nil
	idString := id.String()
	ch := make(chan IncomingMessage, 1)
	c.inFlight[idString] = ch

	return idString, ch
}

func (c *Client) releaseNonce(nonce string) {

	if _, ok := c.inFlight[nonce]; ok == false {
		return
	}

	close(c.inFlight[nonce])
	delete(c.inFlight, nonce)
}

func (c *Client) Request(payload *OutgoingMessage) (ResultFunction, error) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	nonce, ch := c.generateNonce()
	payload.Nonce = nonce

	return func() *IncomingMessage {
		defer c.releaseNonce(nonce)

		for result := range ch {
			return &result
		}

		return nil
	}, c.conn.WriteJSON(&payload)
}

func (c *Client) Close() error {
	atomic.StoreInt64(&c.connectionStatus, 0) //ping loop will die
	return c.conn.Close()
}

/**
Handler triggers for all messages received by the client
*/
func (c *Client) SetCatchAllHandler(h HandlerFunction) {
	c.catchAllHandler = h
}

/**
Handles unknown messages
*/
func (c *Client) SetUnknownHandler(h HandlerFunction) {
	c.unknownHandler = h
}
