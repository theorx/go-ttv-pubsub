package TTVClient

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/theorx/go-ttv-pubsub/pkg/Topic"
	"sync"
	"sync/atomic"
	"time"
)

const TwitchPubSubHost = "wss://pubsub-edge.twitch.tv/"

type Client struct {
	conn *websocket.Conn
	/* Websocket endpoint used for connection */
	endpoint string
	/* Twitch pubsub oauth access_token */
	authToken string
	/* Map on requests in flight with response channels tracked by nonce values*/
	inFlight map[string]chan IncomingMessage
	/* Only accessed via atomic package */
	connectionStatus int64
	/* Slice of topics are tracked to allow automatic resubscribing during reconnect*/
	topics []Topic.Topic
	/* Support for custom loggers to remove log package dependency */
	logFunction LogFunction
	//connection error
	err error
	/* Concurrent writes to websocket client are not supported*/
	writeLock *sync.Mutex
	/* Ping&Pong timestamps to track connection health */
	lastPingTS int64
	lastPongTS int64

	/* Handlers for all supported topics */
	bitsHandler          BitsHandlerFunction
	bitsBadgeHandler     BitsBadgeHandlerFunction
	subscriptionsHandler SubscriptionsHandlerFunction
	commerceHandler      CommerceHandlerFunction
	whispersHandler      WhispersHandlerFunction
	moderationHandler    ModerationActionHandlerFunction
	/* catch all handler will always be triggered no matter the message type */
	catchAllHandler HandlerFunction
	/* unknown handler is only triggered if all of the other handlers fail to handle the message */
	unknownHandler HandlerFunction
}

/*

Creates a new twitch pubsub connection
it is important to call a successful Subscribe() command within 15 seconds
from creating the connection otherwise twitch will close the connection from the server

*/
func CreateClient(authToken string, endpoint string) (*Client, error) {
	client := &Client{
		endpoint:   endpoint,
		authToken:  authToken,
		writeLock:  &sync.Mutex{},
		lastPingTS: time.Now().Unix(),
		lastPongTS: time.Now().Unix(),
		inFlight:   make(map[string]chan IncomingMessage),
		logFunction: func(i ...interface{}) {
			//default null-logger
		},
	}

	return client, client.init()
}

func (c *Client) init() error {
	c.log("Initializing websocket connection..")
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

	c.log("Reconnecting..")
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
			c.log("Failed closing old connection", err)
		}
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.endpoint, nil)
	c.conn = conn

	if err != nil {
		c.err = err
		return
	}
	atomic.StoreInt64(&c.connectionStatus, 1) //connection back up

	//reset the timestamps to current timestamp
	c.lastPongTS = time.Now().Unix()
	c.lastPingTS = time.Now().Unix()

	go c.readLoop()
	go c.pingLoop()
	c.ping()

	//call listen with the stored topics -- first back up the values and then save them after the Subscribe() call to prevent having duplicate values
	if len(c.topics) > 0 {
		c.err = c.Subscribe(c.topics)
	}
	c.log("Reconnecting was successful")
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

func (c *Client) request(payload *OutgoingMessage) (ResultFunction, error) {
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

func (c *Client) log(v ...interface{}) {
	c.logFunction(v...)
}
