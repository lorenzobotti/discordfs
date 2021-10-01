package discordfs

import (
	"encoding/json"
	"io/fs"
	"path"
	"time"
)

type FileInfo struct {
	name       string
	isFolder   bool
	pubblished time.Time
	size       int
}

type PartInfo struct {
	part   int
	length int
	of     int
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

// field names `Name` and `Size` are necessary for fs.FileInfo, se we make them unexported
func (fi FileInfo) MarshalJSON() ([]byte, error) {
	exported := struct {
		Name       string    `json:"name"`
		Pubblished time.Time `json:"pubblished"`
		Size       int       `json:"size"`
	}{fi.name, fi.pubblished, fi.size}
	return json.Marshal(exported)
}

func (fi *FileInfo) UnmarshalJSON(data []byte) error {
	exported := struct {
		Name       string    `json:"name"`
		Pubblished time.Time `json:"pubblished"`
		Size       int       `json:"size"`
	}{}
	err := json.Unmarshal(data, &exported)
	if err != nil {
		return err
	}

	fi.name = exported.Name
	fi.pubblished = exported.Pubblished
	fi.size = exported.Size
	return nil
}

func (pi PartInfo) MarshalJSON() ([]byte, error) {
	exported := struct {
		Part   int `json:"part"`
		Length int `json:"length"`
		Of     int `json:"of"`
	}{pi.part, pi.length, pi.of}
	return json.Marshal(exported)
}

func (pi *PartInfo) UnmarshalJSON(data []byte) error {
	exported := struct {
		Part   int `json:"part"`
		Length int `json:"length"`
		Of     int `json:"of"`
	}{}
	err := json.Unmarshal(data, &exported)
	if err != nil {
		return err
	}

	pi.part = exported.Part
	pi.length = exported.Length
	pi.of = exported.Of
	return nil
}

// implementing fs.FileInfo

func (f FileInfo) Name() string       { return path.Base(f.name) }
func (f FileInfo) Size() int64        { return int64(f.size) }
func (f FileInfo) Mode() fs.FileMode  { return fs.FileMode(0444) }
func (f FileInfo) ModTime() time.Time { return f.pubblished }
func (f FileInfo) IsDir() bool        { return f.isFolder }
func (f FileInfo) Sys() interface{}   { return nil }

func (f FileInfo) FullPath() string { return f.name }

func NewFileInfo(filePath string, pubblished time.Time, size int) FileInfo {
	return FileInfo{
		name:       CleanPath(filePath),
		pubblished: pubblished,
		size:       size,
	}
}
