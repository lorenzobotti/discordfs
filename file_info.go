package discordfs

import (
	"time"
)

type FileInfo struct {
	Name       string    `json:"name"`
	Pubblished time.Time `json:"pubblished"`
}

type PartInfo struct {
	Part   int `json:"part"`
	Length int `json:"length"`
	Of     int `json:"of"`
}

type ChunkInfo struct {
	File FileInfo `json:"file"`
	Part PartInfo `json:"part"`
	Url  string   `json:"-"`
}

type FileChunk struct {
	Info     ChunkInfo
	Contents []byte
}
