package discordfs

import (
	"io"
)

type chunker struct {
	info FileInfo
	r    io.Reader

	lastChunk  int
	chunkSize  int
	chunksRead int
}

func newChunker(info FileInfo, file io.Reader, chunkSize int) chunker {
	return chunker{
		info:       info,
		r:          file,
		lastChunk:  chunksNeeded(info.size, chunkSize) - 1,
		chunkSize:  chunkSize,
		chunksRead: 0,
	}
}

// nextChunk returns the next piece of the file. if it's the last piece
// `done` is set to true
func (c *chunker) nextChunk() (chunk FileChunk, done bool, err error) {
	//var n int
	//n, err = c.r.Read(contents)
	//if n < c.chunkSize {
	//	contents = contents[:n]
	//}

	if c.chunksRead > c.lastChunk {
		done = true
		return
	}

	howMuchToRead := c.chunkSize
	if c.chunksRead == c.lastChunk {
		partialChunkLeft := c.info.size % c.chunkSize
		if partialChunkLeft > 0 {
			howMuchToRead = partialChunkLeft
		}
	} else if c.chunksRead > c.lastChunk {
		return
	}

	buf := make([]byte, howMuchToRead)

	_, err = insistReading(c.r, buf)
	if err != nil {
		return
	}

	chunk = FileChunk{
		Contents: buf,
		Info: ChunkInfo{
			File: c.info,
			Part: PartInfo{
				part:   c.chunksRead,
				length: len(buf),
			},
		},
	}

	c.chunksRead += 1

	return
}

type insistentReader struct{ r io.Reader }

func (ir insistentReader) Read(dt []byte) (int, error) {
	sliceLeft := dt
	howMuchToRead := len(dt)
	read := 0

	for {
		n, err := ir.r.Read(sliceLeft)
		read += n

		if n < len(sliceLeft) {
			return read, err
		}

		if read == howMuchToRead {
			return howMuchToRead, nil
		} else if read < howMuchToRead {
			sliceLeft = dt[:(howMuchToRead - read)]
		}
	}
}

func insistReading(r io.Reader, dt []byte) (int, error) {
	return insistentReader{r}.Read(dt)
}
