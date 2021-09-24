package discordfs

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

// questa funzione gli passo un url e un chunk da riempire
// (tramite puntatore). mi dirai, non è meglio fargli riempire
// solo la slice di byte, che è più generico quindi più riutilizzabile?
// ottima domanda, la soluzione è che è complicato, te
// lo spiego a casa se vuoi
func downloadIntoChunk(url string, into *FileChunk) {
	req, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, req.Body)
	if err != nil {
		panic(err)
	}

	into.Contents = buf.Bytes()
}

// stesso di downloadIntoChunk ma segnala al WaitGroup quando ha finito
func downloadIntoChunkWG(url string, into *FileChunk, wg *sync.WaitGroup) {
	downloadIntoChunk(url, into)
	wg.Done()
}

func downloadIntoChunkChan(url string, into *FileChunk, done chan<- bool) {
	downloadIntoChunk(url, into)
	done <- true
}
