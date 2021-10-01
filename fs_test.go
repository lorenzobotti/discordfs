package discordfs

import (
	"io/fs"
	"testing"
)

func TestImplementsFS(t *testing.T) {
	_ = (fs.FS)(DiscStorage{})
}

func TestImplementsStatFS(t *testing.T) {
	_ = (fs.StatFS)(DiscStorage{})
}

func TestImplementsReadDirFS(t *testing.T) {
	_ = (fs.ReadDirFS)(DiscStorage{})
}
