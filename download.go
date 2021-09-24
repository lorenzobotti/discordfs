package discordfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	dg "github.com/bwmarrin/discordgo"
)

// Receive looks for the file named with the provided name
// in the channel with the provided channelID with the provided
// session, storing the result in the provided Writer
func Receive(into io.Writer, s *dg.Session, channelID, name string) error {
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
				fmt.Println("no attachments found")
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
	for i, chunk := range pieces {
		fmt.Println("starting download of chunk", i)
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
