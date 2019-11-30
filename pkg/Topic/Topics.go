package Topic

import "fmt"

/**
https://dev.twitch.tv/docs/pubsub

*/
type Topic string

func Bits(channelID int) Topic {
	return Topic(fmt.Sprintf("channel-bits-events-v2.%d", channelID))
}

func BitsBadgeNotification(channelID int) Topic {
	return Topic(fmt.Sprintf("channel-bits-badge-unlocks.%d", channelID))
}

func Subscriptions(channelID int) Topic {
	return Topic(fmt.Sprintf("channel-subscribe-events-v1.%d", channelID))
}

func Commerce(channelID int) Topic {
	return Topic(fmt.Sprintf("channel-commerce-events-v1.%d", channelID))
}

func Whispers(channelID int) Topic {
	return Topic(fmt.Sprintf("whispers.%d", channelID))
}

func ModerationAction(userID int, channelID int) Topic {
	return Topic(fmt.Sprintf("chat_moderator_actions.%d.%d", userID, channelID))
}
