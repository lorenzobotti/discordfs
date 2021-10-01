package discordfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const Root = "/"

// ListFiles lists all the files `st` can find
func (st DiscStorage) ListFiles(folder string) ([]FileInfo, error) {
	folder = CleanPath(folder)

	iter := newMessageIterator(st.session, st.channelId)
	set := map[string]FileInfo{}

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

		if strings.HasPrefix(info.File.name, folder) {
			relativePath := strings.TrimPrefix(info.File.name, folder)

			if !isInFolder(relativePath) {
				set[info.File.name] = info.File
			} else {
				folderName := topFolder(info.File.name)
				set[folderName] = FileInfo{
					name:     folderName,
					isFolder: true,
				}
			}
		}
	}

	output := make([]FileInfo, 0, len(set))
	for _, info := range set {
		output = append(output, info)
	}

	return output, nil
}

var ErrFileNotFound = errors.New("file couldn't be found")

// Open returns a DiscFile if found on the channel
func (st DiscStorage) GetFile(filename string) (*DiscFile, error) {
	filename = CleanPath(filename)
	iter := newMessageIterator(st.session, st.channelId)

	for {
		msg, err := iter.next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("errore nel richiedere messaggi: %w", err)
		}

		info := ChunkInfo{}
		err = json.Unmarshal([]byte(msg.Content), &info)
		if err != nil {
			continue
		}

		if info.File.name == filename {
			// todo: newDiscFile()
			return &DiscFile{
				storage: st,
				info:    info.File,
			}, nil
		}
	}

	return nil, ErrFileNotFound
}

// DoesFileExist checks if a file exists on the cloud channel
func (st DiscStorage) DoesFileExist(filename string) (bool, error) {
	filename = CleanPath(filename)
	_, err := st.Open(filename)

	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}
