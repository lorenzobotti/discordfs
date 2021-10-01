package discordfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
)

// ReceiveAllAtOnce looks for the file named `name`
// in the channel with the provided `channelID` with the provided
// session, storing the result into `into`. It first downloads
// the whole file, then writes it all at once into `into`.
// If you'd rather get each piece as soon as it arrives try using Receive()
func (st DiscStorage) ReceiveAllAtOnce(into io.Writer, filename string) error {
	filename = CleanPath(filename)
	chunks, err := st.fileChunks(filename)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(chunks))
	for _, chunk := range chunks {
		go downloadIntoChunkWG(chunk.Info.Url, chunk, wg)
	}

	// waiting for all the downloads to finish
	wg.Wait()

	// writing each chunk in order
	for _, chunk := range chunks {
		_, err := into.Write(chunk.Contents)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReceiveAllAtOnce looks for the file named `name`
// in the channel with the provided `channelID` with the provided
// session, storing the result into `into`. It writes each piece into
// `into` as soon as it comes
// todo: remove duplication between this and ReceiveAllAtOnce()?
func (st DiscStorage) Receive(into io.Writer, filename string) error {
	filename = CleanPath(filename)
	chunks, err := st.fileChunks(filename)
	if err != nil {
		return err
	}

	howManyChunks := len(chunks)

	chunkChannels := make([]chan bool, howManyChunks)
	for i, chunk := range chunks {
		chunkChannels[i] = make(chan bool, 1)
		go downloadIntoChunkChan(chunk.Info.Url, chunk, chunkChannels[i])
	}

	for i := 0; i < len(chunks); i++ {
		chunk := chunks[i]
		if chunk == nil {
			continue
		}
		// we wait for the next chunk to be ready
		// by doing this we can send the chunk's contents on the
		// writer as soon as we download it, laying the foundations
		// for some streaming capability in the future, maybe
		<-chunkChannels[i]

		emptyChunk := &FileChunk{}
		// checking if the download came through correctly
		if chunk == emptyChunk {
			return fmt.Errorf("missing chunk %d", i)
		}

		_, err := into.Write(chunk.Contents)
		if err != nil {
			return err
		}
	}

	return nil
}

// fileChunks looks in the channel for all chunks of the file, unmarshals them and makes sure
// they're all there
func (st DiscStorage) fileChunks(filename string) ([]*FileChunk, error) {
	iter := newMessageIterator(st.session, st.channelId)

	// keeps the chunks found
	pieces := []*FileChunk{}

	piecesFound := 0
	howManyChunks := 0

	// cerco, tra tutti i messaggi, quelli che hanno il file che mi serve
	// e mi salvo l'url, non scarico niente ancora
	for {
		mess, err := iter.next()
		if err != nil {
			// if it's gone through all the messages in the channel
			if err == io.EOF {
				break
				// ...or if it's some other kind of error
			} else {
				return pieces, err
			}
		}

		// if it can't read the message as json it just assumes it's not relevant
		info := ChunkInfo{}
		err = json.Unmarshal([]byte(mess.Content), &info)
		if err != nil {
			continue
		}

		// if i've found the file i'm looking for
		if info.File.name == filename {
			if len(mess.Attachments) == 0 {
				// todo: should this be some kind of error?
				continue
			}

			thisIsTheFirstChunk := howManyChunks == 0
			if thisIsTheFirstChunk {
				howManyChunks = info.Part.of + 1
				pieces = make([]*FileChunk, howManyChunks)
			}

			info.Url = mess.Attachments[0].URL
			pieces[info.Part.part] = &FileChunk{
				Info: info,
			}
			piecesFound += 1

			// se ho trovato tutti i pezzi
			if piecesFound == howManyChunks {
				break
			}
		}
	}

	if howManyChunks == 0 {
		return nil, errors.New("file couldn't be found")
	}

	if piecesFound != howManyChunks {
		return nil, errors.New("couldn't find all chunks")
	}

	return pieces, nil
}
