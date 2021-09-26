package discordfs

import (
	"io/fs"
	"testing"
)

func TestImplementsFS(t *testing.T) {
	st, err := newTestStorage()
	if err != nil {
		t.Fatalf("error initializing storage: %s", err.Error())
	}

	functionThatTakesFS := func(_ fs.FS) {}
	functionThatTakesFS(st)
}
