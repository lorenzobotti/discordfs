package discordfs

import (
	"encoding/json"
	"errors"

	dg "github.com/bwmarrin/discordgo"
)

func getChannels(s *dg.Session) ([]*dg.Channel, error) {
	err := checkWebSocket(s)
	if err != nil {
		return nil, err
	}

	output := []*dg.Channel{}
	for _, guild := range s.State.Ready.Guilds {
		channels, err := s.GuildChannels(guild.ID)
		if err != nil {
			return nil, err
		}

		for _, channel := range channels {
			if channel.Type != dg.ChannelTypeGuildText {
				continue
			}

			output = append(output, channel)
		}
	}

	return output, nil
}

// GetCloudChannel looks through all the channels the bot is in and
// returns the first one it thinks to be a valid cloud channel
func GetCloudChannel(s *dg.Session) (*dg.Channel, error) {
	channels, err := getChannels(s)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		isCloud, err := IsCloudChannel(s, channel.ID)
		if isCloud {
			return channel, err
		}

		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("no cloud channel found")
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

// connects to the websocket if it hasn't already
func checkWebSocket(s *dg.Session) error {
	err := s.Open()
	if err != nil && err != dg.ErrWSAlreadyOpen {
		return err
	} else {
		return nil
	}
}
