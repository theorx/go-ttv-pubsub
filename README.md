# Twitch.tv pubsub client

## Author
* Lauri Orgla `theorx@hotmail.com`

## Description
*PUBSUB Websocket api client implementation written in golang, motivation was to build a client
that actually works. There are other implementations that are not working so well or are not 
working at all. This implementation covers all publicly documented behaviours of the pubsub api
and also one undocumented `moderation actions` topic.*

* Documentation to official twitch api [Twitch pubsub doc](https://dev.twitch.tv/docs/pubsub)

## Connection handling
*Twitch requires to send a __PING__ message at least once per 5 minutes. This implementation is
sending a __PING__ every 60 seconds.*

*Twitch states that every ping request __should receive a response within 10 seconds__, this 
implementation will wait __up to 20 seconds__ before going to a __reconnect state__*

*In case of __RECONNECT__ or failing to receive a __PONG__ all the current topics
that you have subscribed to __will get resubscribed automatically__*

__For example:__

```go

//you create a client c
c.Subscribe(..)
//now RECONNECT happens
//you will not need to call Subscribe again, the topics will be subscribed to
//on the new connection automatically

```

## API Overview

* Set\*Handler functions are for setting handler functions to the types of topics you have subscribed to

```go

	client.SetCatchAllHandler(func(message TTVClient.IncomingMessage) {})
	client.SetLogFunction(func(i ...interface{}) {})
	client.SetBitsHandler(func(message TTVClient.BitsMsg) {})
	client.SetModerationHandler(func(message TTVClient.ModerationActionMsg) {})
	client.SetCommerceHandler(func(message TTVClient.CommerceMsg) {})
	client.SetBitsBadgeHandler(func(message TTVClient.BitsBadgeMsg) {})
	client.SetSubscriptionsHandler(func(message TTVClient.SubscriptionMsg) {})
	client.SetWhisperHandler(func(message TTVClient.WhisperMsg) {})
	client.SetUnknownHandler(func(message TTVClient.IncomingMessage) {})
```

* CreateClient creates the client, expects access token and twitch websocket endpoint
* There is one predefined constant: __TTVClient.TwitchPubSubHost__ which corresponds to: wss://pubsub-edge.twitch.tv/
* Endpoint required for __CreateClient__ requires correct prefix `wss://`
```go

	TTVClient.CreateClient("your-access_token", TTVClient.TwitchPubSubHost)

```

* Subscribe/unsubscribe have the same signature
* Both of the operations require a slice of Topics - `[]Topic.Topic`

```go
	client.Subscribe([]Topic.Topic{
		Topic.Bits(1),
		Topic.BitsBadgeNotification(1),
		Topic.Commerce(1),
		Topic.Whispers(1),
		Topic.Subscriptions(1),
		Topic.ModerationAction(1, 2),
	})

	client.Unsubscribe([]Topic.Topic{
		Topic.Bits(1),
		Topic.BitsBadgeNotification(1),
		Topic.Commerce(1),
		Topic.Whispers(1),
		Topic.Subscriptions(1),
		Topic.ModerationAction(1, 2),
	})

```

* Topic.Bits(`channel id`)
* Topic.BitsBadgeNotification(`channel id`)
* Topic.Commerce(`channel id`)
* Topic.Whispers(`user id`)
* Topic.Subscriptions(`channel id`)
* Topic.ModerationAction(`user id`, `channel id`)

* Closing the connection - `client.Close()`

### Errors

* TTVClient.ErrorOperationFailed = errors.New("sub/unsub operation failed") <- In this case usually the credentials are incorrect
* TTVClient.ErrorNotConnected = errors.New("not connected") <- Connection is down / TTVClient has been closed


## Usage example

```go

package main

import (
	"github.com/theorx/go-ttv-pubsub/pkg/TTVClient"
	"github.com/theorx/go-ttv-pubsub/pkg/Topic"
	"log"
	"time"
)

func main() {

	//create a connection
	client, err := TTVClient.CreateClient("your-access_token", TTVClient.TwitchPubSubHost)

	if err != nil {
		log.Fatal("Failed to connect", err)
		return
	}

	//set up handlers before subscribing

	client.SetModerationHandler(func(message TTVClient.ModerationActionMsg) {
		log.Println("Moderation event received", message)
	})

	err = client.Subscribe(
		[]Topic.Topic{
			Topic.ModerationAction(64417816, 64417816),
		},
	)

	if err != nil {
		log.Println(err)
		return
	}

	time.Sleep(time.Second * 60)

	log.Println(client.Close())
}

```