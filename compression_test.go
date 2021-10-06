package discordfs

import (
	"bytes"
	"crypto/sha256"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func TestSendCompressed(t *testing.T) {
	st, err := newTestStorage()
	if err != nil {
		t.Fatalf("can't open storage: %s", err.Error())
	}

	err = deleteTestFiles(st)
	if err != nil {
		t.Fatalf("can't delete test files: %s", err.Error())
	}

	fileContents, err := ioutil.ReadFile(filepath.Join(testFilesDir, "pecunia non olet.pdf"))
	if err != nil {
		t.Fatalf("Error reading file(): %s", err.Error())
	}

	originalSum := sha256.Sum256(fileContents)
	file := bytes.NewBuffer(fileContents)

	err = st.SendCompressed(
		file,
		NewFileInfo("test_file_compressed", time.Now(), len(fileContents)),
		4*MB,
	)
	if err != nil {
		t.Fatalf("Error in SendCompressed(): %s", err.Error())
	}

	downloadedBuf := bytes.Buffer{}
	err = st.ReceiveCompressed(&downloadedBuf, "test_file_compressed")
	if err != nil {
		t.Fatalf("Error in ReceiveCompressed(): %s", err.Error())
	}

	downloadedSum := sha256.Sum256(downloadedBuf.Bytes())

	if originalSum != downloadedSum {
		t.Fatal("original and downloaded checksum don't match")
	}
}
