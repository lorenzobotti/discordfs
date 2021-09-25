package discordfs

import (
	"encoding/json"
	"io"
)

func (st DiscStorage) Delete(filename string) error {
	iter := newMessageIterator(st.session, st.channelId)

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
			err = st.session.ChannelMessageDelete(st.channelId, mess.ID)
			if err != nil {
				return err
			}
		}
	}
}
