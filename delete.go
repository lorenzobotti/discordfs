package discordfs

import (
	"encoding/json"
	"io"

	dg "github.com/bwmarrin/discordgo"
)

func Delete(s *dg.Session, channelID, filename string) error {
	iter := newMessageIterator(s, channelID)

	for {
		mess, err := iter.next()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return err
			}
		}

		// todo: make this its own function
		info := ChunkInfo{}
		err = json.Unmarshal([]byte(mess.Content), &info)
		if err != nil {
			continue
		}

		if info.File.Name == filename {
			err = s.ChannelMessageDelete(channelID, mess.ID)
			if err != nil {
				return err
			}
		}
	}
}
