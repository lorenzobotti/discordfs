package discordfs

import (
	"encoding/json"
	"fmt"
	"io"
)

// ListFiles lists all the files `st` can find
// todo: make this return a slice of *os.FileInfo
func (st DiscStorage) ListFiles() ([]string, error) {
	iter := newMessageIterator(st.session, st.channelId)
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

// DoesFileExist checks if a file exists on the cloud channel
func (st DiscStorage) DoesFileExist(filename string) (bool, error) {
	iter := newMessageIterator(st.session, st.channelId)

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
