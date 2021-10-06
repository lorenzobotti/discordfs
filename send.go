package discordfs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	dg "github.com/bwmarrin/discordgo"
)

// Send splits the `file` into chunks of size `chunkSize` and sends each one
func (st DiscStorage) Send(file io.Reader, info FileInfo, chunkSize int) error {
	info.name = cleanPath(info.name)
	div := newChunker(info, file, chunkSize)

	for {
		chunk, done, err := div.nextChunk()
		if err != nil {
			return err
		}

		if done {
			break
		}

		if len(chunk.Contents) == 0 {
			return errors.New("empty chunk")
		}

		info, err := json.Marshal(chunk.Info)
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("%02d_%s", chunk.Info.Part.part, chunk.Info.File.name)

		attachment := dg.File{
			Name:   filename,
			Reader: bytes.NewBuffer(chunk.Contents),
		}

		message := dg.MessageSend{
			Content: string(info),
			Files:   []*dg.File{&attachment},
		}

		// todo: retry on fail
		_, err = st.session.ChannelMessageSendComplex(st.channelId, &message)
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

	info := FileInfo{
		name:       path.Base(file.Name()),
		pubblished: stat.ModTime(),
		size:       int(stat.Size()),
	}

	return st.Send(file, info, chunkSize)
}

func (st DiscStorage) SendCompressed(file io.Reader, info FileInfo, chunkSize int) error {
	buf := bytes.Buffer{}
	err := readAndCompressInto(&buf, file, BestSpeed)
	if err != nil {
		return fmt.Errorf("error in readAndCompress(): %w", err)
	}

	info.size = buf.Len()
	err = st.Send(&buf, info, chunkSize)
	if err != nil {
		return fmt.Errorf("error in Send(): %w", err)
	}

	return nil
}

func chunksNeeded(fileSize, chunkSize int) int {
	chunks := fileSize / chunkSize

	if chunks*chunkSize >= fileSize {
		return chunks
	} else {
		return chunks + 1
	}
}
