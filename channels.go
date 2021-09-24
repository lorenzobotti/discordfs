package discordfs

import (
	"encoding/json"

	dg "github.com/bwmarrin/discordgo"
)

func getChannels(s *dg.Session) []*dg.Channel {
	channels := []*dg.Channel{}
	for _, guild := range s.State.Guilds {
		for _, channel := range guild.Channels {
			if channel.Type != dg.ChannelTypeGuildText {
				continue
			}

			channels = append(channels, channel)
		}
	}

	return channels
}

// GetCloudChannel looks through all the channels the bot is in and
// returns the first one it thinks to be a valid cloud channel
func GetCloudChannel(s *dg.Session) *dg.Channel {
	channels := getChannels(s)

	for _, channel := range channels {
		isCloud, _ := IsCloudChannel(s, channel.ID)
		if isCloud {
			return channel
		}
	}

	return nil
}

func IsCloudChannel(s *dg.Session, channelID string) (bool, error) {
	messages, err := s.ChannelMessages(channelID, messageLimit, "", "", "")
	if err != nil {
		return false, err
	}

	for _, mess := range messages {
		// controllo che il messaggio sia uno di quelli generati
		// da questo programma
		t := ChunkInfo{}
		err := json.Unmarshal([]byte(mess.Content), &t)
		// se non dà errore è un messaggio valido
		if err == nil {
			return true, nil
		}
	}

	return false, nil
}
