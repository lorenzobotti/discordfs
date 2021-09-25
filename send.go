package discordfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Send splits the `file` into chunks of size `chunkSize` and sends each one
func (st DiscStorage) Send(file io.Reader, filename string, chunkSize, fileSize int) error {
	div := newChunker(FileInfo{
		Name: filename,
		// TODO: take this as input
		Pubblished: time.Now(),
	}, file, chunkSize)

	lastChunk := chunksNeeded(fileSize, chunkSize) - 1

	for {
		chunk, done, err := div.nextChunk()
		chunk.Info.Part.Of = lastChunk
		if err != nil {
			return err
		}

		if done {
			break
		}

		info, err := json.Marshal(chunk.Info)
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("%02d_%s", chunk.Info.Part.Part, chunk.Info.File.Name)

		_, err = st.session.ChannelFileSendWithMessage(st.channelId, string(info), filename, bytes.NewBuffer(chunk.Contents))
		if err != nil {
			return err
		}

	}

	return nil
}

// SendFile is a frontend to `Send` that doesn't ask you for a name or a file size
func (st DiscStorage) SendFile(file *os.File, chunkSize int) error {
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error calling file.Stat(): %w", err)
	}

	return st.Send(file, filepath.Base(file.Name()), chunkSize, int(stat.Size()))
}

func chunksNeeded(fileSize, chunkSize int) int {
	chunks := fileSize / chunkSize

	if chunks*chunkSize >= fileSize {
		return chunks
	} else {
		return chunks + 1
	}
}
