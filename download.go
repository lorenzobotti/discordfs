package discordfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	dg "github.com/bwmarrin/discordgo"
)

// ReceiveAllAtOnce looks for the file named with the provided name
// in the channel with the provided channelID with the provided
// session, storing the result in the provided Writer
func ReceiveAllAtOnce(into io.Writer, s *dg.Session, channelID, name string) error {
	iter := newMessageIterator(s, channelID)

	// dato un singolo file non so quanti siano i pezzi, perciò
	// invece di usare un array uso un dizionario (map) con numeri
	// come chiavi
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

		// per leggere json si fa così purtroppo
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

			// se ho trovato tutti i pezzi
			if len(pieces) == howManyChunks {
				break
			}
		}
	}

	if len(pieces) == 0 {
		return errors.New("file name couldn't be found")
	}

	// un WaitGroup serve a sincronizzare più thread
	// a ogni thread passo un puntatore a un chunk e gli
	// dò la responsabilità di "riempire" quel chunk
	// quando ha finito segnala al WaitGroup che ha finito
	// (tramite wg.Done())
	wg := &sync.WaitGroup{}
	wg.Add(len(pieces))
	for _, chunk := range pieces {
		go downloadIntoChunkWG(chunk.Info.Url, chunk, wg)
	}

	// aspetto finchè non finiscono di scaricare tutti i chunk
	wg.Wait()

	// vado a scrivere ciascun chunk in ordine nel file
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

func Receive(into io.Writer, s *dg.Session, channelID, name string) error {
	iter := newMessageIterator(s, channelID)

	// keeps the chunks found
	pieces := []*FileChunk{}

	// keeps track of how many we found already
	piecesFound := map[int]struct{}{}
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

		// per leggere json si fa così purtroppo
		info := ChunkInfo{}
		err = json.Unmarshal([]byte(mess.Content), &info)
		if err != nil {
			continue
		}

		if info.File.Name == name {
			if len(mess.Attachments) == 0 {
				fmt.Println("no attachments found")
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
		chunkChannels[i] = make(chan bool)
		go downloadIntoChunkChan(chunk.Info.Url, chunk, chunkChannels[i])
	}

	// vado a scrivere ciascun chunk in ordine nel file
	for i := 0; i < len(pieces); i++ {
		// we wait for the right chunk to be ready
		<-chunkChannels[i]

		chunk := pieces[i]
		emptyChunk := &FileChunk{}
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
