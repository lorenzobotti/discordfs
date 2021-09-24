package discordfs

import (
	"bytes"
	"fmt"
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
	fmt.Println("i start to download")
	req, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	fmt.Println("finished downloading")

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, req.Body)
	fmt.Println("copied into buffer")
	if err != nil {
		panic(err)
	}

	fmt.Println("filling chunk")
	into.Contents = buf.Bytes()
}

// stesso di downloadIntoChunk ma segnala al WaitGroup quando ha finito
func downloadIntoChunkWG(url string, into *FileChunk, wg *sync.WaitGroup) {
	downloadIntoChunk(url, into)
	wg.Done()
}
