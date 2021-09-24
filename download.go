package discordfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	dg "github.com/bwmarrin/discordgo"
)

// ReceiveAllAtOnce looks for the file named `name`
// in the channel with the provided `channelID` with the provided
// session, storing the result into `into`. It first downloads
// the whole file, then writes it all at once into `into`.
// If you'd rather get each piece as soon as it arrives try using Receive()
func ReceiveAllAtOnce(into io.Writer, s *dg.Session, channelID, name string) error {
	iter := newMessageIterator(s, channelID)

	// todo: make this not need a map anymore
	// (it used to before i implemented the FilePart.Of field)
	pieces := map[int]*FileChunk{}
	howManyChunks := 0

	// cerco, tra tutti i messaggi, quelli che hanno il file che mi serve
	// e mi salvo l'url, non scarico niente ancora
	for {
		mess, err := iter.next()
		if err != nil {
			// se ho finito i messaggi nel canale (EOF sta per End Of File)
			if err == io.EOF {
				break
				// ...altri tipi di errore
			} else {
				return err
			}
		}

		info := ChunkInfo{}
		err = json.Unmarshal([]byte(mess.Content), &info)
		if err != nil {
			continue
		}

		if info.File.Name == name {
			if len(mess.Attachments) == 0 {
				continue
			}

			info.Url = mess.Attachments[0].URL
			howManyChunks = info.Part.Of + 1

			pieces[info.Part.Part] = &FileChunk{
				Info: info,
			}

			// if i found all the pieces
			if len(pieces) == howManyChunks {
				break
			}
		}
	}

	if len(pieces) == 0 {
		return errors.New("file name couldn't be found")
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(pieces))
	for _, chunk := range pieces {
		go downloadIntoChunkWG(chunk.Info.Url, chunk, wg)
	}

	// waiting for all the downloads to finish
	wg.Wait()

	// writing each chunk in order
	for i := 0; i < len(pieces); i++ {
		chunk, ok := pieces[i]
		if !ok {
			return fmt.Errorf("missing chunk %d", i)
		}

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
func Receive(into io.Writer, s *dg.Session, channelID, name string) error {
	iter := newMessageIterator(s, channelID)

	// keeps the chunks found
	pieces := []*FileChunk{}

	// keeps track of how many we found already
	// it sorta behaves like a set
	piecesFound := map[int]struct{}{}
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
				return err
			}
		}

		// if it can't read the message as json it just assumes it's not relevant
		info := ChunkInfo{}
		err = json.Unmarshal([]byte(mess.Content), &info)
		if err != nil {
			continue
		}

		// if i've found the file i'm looking for
		if info.File.Name == name {
			if len(mess.Attachments) == 0 {
				// todo: should this be some kind of error?
				continue
			}

			thisIsTheFirstChunk := howManyChunks == 0
			if thisIsTheFirstChunk {
				howManyChunks = info.Part.Of + 1
				pieces = make([]*FileChunk, howManyChunks)
			}

			info.Url = mess.Attachments[0].URL
			pieces[info.Part.Part] = &FileChunk{
				Info: info,
			}
			piecesFound[info.Part.Part] = struct{}{}

			// se ho trovato tutti i pezzi
			if len(piecesFound) == howManyChunks {
				break
			}
		}
	}

	if len(pieces) == 0 {
		return errors.New("file name couldn't be found")
	}

	chunkChannels := make([]chan bool, howManyChunks)
	for i, chunk := range pieces {
		chunkChannels[i] = make(chan bool, 1)
		go downloadIntoChunkChan(chunk.Info.Url, chunk, chunkChannels[i])
	}

	for i := 0; i < len(pieces); i++ {
		// we wait for the next chunk to be ready
		// by doing this we can send the chunk's contents on the
		// writer as soon as we download it, laying the foundations
		// for some streaming capability in the future, maybe
		<-chunkChannels[i]

		chunk := pieces[i]
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
