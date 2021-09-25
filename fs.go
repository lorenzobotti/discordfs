package discordfs

import "github.com/bwmarrin/discordgo"

// DiscStorage contains the necessary information to connect to the cloud channel
type DiscStorage struct {
	session   *discordgo.Session
	channelId string
}

// NewStorage builds a new DiscStorage
func NewStorage(s *discordgo.Session, channelId string) DiscStorage {
	return DiscStorage{
		session:   s,
		channelId: channelId,
	}
}
