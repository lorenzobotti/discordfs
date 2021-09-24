package discordfs

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/bwmarrin/discordgo"
)

func ListFiles(s *discordgo.Session, channelId string) ([]string, error) {
	iter := newMessageIterator(s, channelId)
	set := map[string]struct{}{}

	for {
		msg, err := iter.next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("can't load next message: %w", err)
		}

		info := ChunkInfo{}
		err = json.Unmarshal([]byte(msg.Content), &info)
		if err != nil {
			continue
		}

		set[info.File.Name] = struct{}{}
	}

	output := make([]string, 0, len(set))
	for name, _ := range set {
		output = append(output, name)
	}

	return output, nil
}

func DoesFileExist(s *discordgo.Session, channelId, filename string) (bool, error) {
	iter := newMessageIterator(s, channelId)

	for {
		msg, err := iter.next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return false, fmt.Errorf("errore nel richiedere messaggi: %w", err)
		}

		info := ChunkInfo{}
		err = json.Unmarshal([]byte(msg.Content), &info)
		if err != nil {
			continue
		}

		if info.File.Name == filename {
			return true, nil
		}
	}

	return false, nil
}
