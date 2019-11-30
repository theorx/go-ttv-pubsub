package WSClient

type HandlerFunction func(message IncomingMessage)
type BitsHandlerFunction func(message BitsMsg)
type BitsBadgeHandlerFunction func(message BitsBadgeMsg)
type SubscriptionsHandlerFunction func(message SubscriptionMsg)
type CommerceHandlerFunction func(message CommerceMsg)
type WhispersHandlerFunction func(message WhisperMsg)
type ModerationActionHandlerFunction func(message ModerationActionMsg)
type ResultFunction func() *IncomingMessage

func (c *Client) SetBitsHandler(h BitsHandlerFunction) {
	c.bitsHandler = h
}

func (c *Client) SetBitsBadgeHandler(h BitsBadgeHandlerFunction) {
	c.bitsBadgeHandler = h
}

func (c *Client) SetSubscriptionsHandler(h SubscriptionsHandlerFunction) {
	c.subscriptionsHandler = h
}

func (c *Client) SetCommerceHandler(h CommerceHandlerFunction) {
	c.commerceHandler = h
}

func (c *Client) SetWhisperHandler(h WhispersHandlerFunction) {
	c.whispersHandler = h
}

func (c *Client) SetModerationHandler(h ModerationActionHandlerFunction) {
	c.moderationHandler = h
}
