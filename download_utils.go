package discordfs

import (
	"bytes"
	"io"
	"net/http"
	"sync"
)

// downloadIntoChunk takes a pointer to a FileChunk and "fills it in" with the resource
// at the specified url
// todo: make this return the error instead of panicking
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

// downloadIntoChunkWG is the same as downloadIntoChunk but decrements the WaitGroup when it's done
func downloadIntoChunkWG(url string, into *FileChunk, wg *sync.WaitGroup) {
	downloadIntoChunk(url, into)
	wg.Done()
}

// downloadIntoChunkChan is the same as downloadIntoChunk but sends on the given channel when it's done
func downloadIntoChunkChan(url string, into *FileChunk, done chan<- bool) {
	downloadIntoChunk(url, into)
	done <- true
}
