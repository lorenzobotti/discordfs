package discordfs

import "io"

type chunker struct {
	info       FileInfo
	r          io.Reader
	chunkSize  int
	chunksRead int
}

func newChunker(info FileInfo, file io.Reader, chunkSize int) chunker {
	return chunker{
		info:       info,
		r:          file,
		chunkSize:  chunkSize,
		chunksRead: 0,
	}
}

// nextChunk returns the next piece of the file. if it's the last piece
// `done` is set to true
func (c *chunker) nextChunk() (chunk FileChunk, done bool, err error) {
	contents := make([]byte, c.chunkSize)

	var n int
	n, err = c.r.Read(contents)
	if n < c.chunkSize {
		contents = contents[:n]
	}

	if err == io.EOF {
		err = nil
		done = true
		return
	}

	chunk = FileChunk{
		Contents: contents,
		Info: ChunkInfo{
			File: c.info,
			Part: PartInfo{
				Part:   c.chunksRead,
				Length: c.chunkSize,
			},
		},
	}
	c.chunksRead += 1
	done = false

	return
}
