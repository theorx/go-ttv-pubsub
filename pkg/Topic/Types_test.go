package Topic_test

import (
	"github.com/theorx/go-ttv-pubsub/pkg/Topic"
	"testing"
)

func TestGetType(t *testing.T) {

	type TestCase struct {
		Type  Topic.Type
		Input string
	}

	table := []TestCase{
		{Topic.TypeInvalid, "test"},
		{Topic.TypeInvalid, ""},
		{Topic.TypeBits, "channel-bits-events-v2.46024993"},
		{Topic.TypeWhispers, "whispers.44322889"},
		{Topic.TypeCommerce, "channel-commerce-events-v1.44322889"},
		{Topic.TypeSubscriptions, "channel-subscribe-events-v1.44322889"},
		{Topic.TypeBitsBadgeNotification, "channel-bits-badge-unlocks.44322889"},
		{Topic.TypeModerationAction, "chat_moderator_actions.test.test"},
	}

	for _, test := range table {
		result := Topic.GetType(test.Input)
		if result != test.Type {
			t.Errorf("Failed getting type expected '%v' but got '%v'", result, test.Type)
		}
	}
}
