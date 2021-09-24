package discordfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Send splits the `file` into chunks of size `chunkSize` and sends each one
func Send(file io.Reader, s *discordgo.Session, channelID, filename string, chunkSize, fileSize int) error {
	div := newChunker(FileInfo{
		Name: filename,
		// TODO: take this as input
		Pubblished: time.Now(),
	}, file, chunkSize)

	lastChunk := fileSize / chunkSize

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

		_, err = s.ChannelFileSendWithMessage(channelID, string(info), filename, bytes.NewBuffer(chunk.Contents))
		if err != nil {
			return err
		}

	}

	return nil
}

// SendFile is a frontend to `Send` that doesn't ask you for a name or a file size
func SendFile(file *os.File, s *discordgo.Session, channelID string, chunkSize int) error {
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error calling file.Stat(): %w", err)
	}

	return Send(file, s, channelID, filepath.Base(file.Name()), chunkSize, int(stat.Size()))
}
