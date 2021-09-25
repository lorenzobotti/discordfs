package discordfs

import "testing"

func TestGetCloudChannel(t *testing.T) {
	sess, err := newTestSession()
	if err != nil {
		t.Fatalf("error connecting to discord: %s", err.Error())
	}

	channel, err := GetCloudChannel(sess)
	if err != nil {
		// todo: provide more info
		t.Fail()
	}
	if channel == nil {
		t.Fatalf("GetCloudChannel() returned nil")
	}

	if channel.ID != channelId {
		t.Fatalf(
			"GetCloudChannel() found the wrong channel: expected %s, found %s",
			channelId,
			channel.ID,
		)
	}
}
