package discordfs

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"io"
	"log"
	"os"
	"path"
	"testing"

	dg "github.com/bwmarrin/discordgo"
)

const maxDownloadSize = 20 * MB
const limitDownloadSize = true
const testFilesDir = "test_files"

func TestNewDownload(t *testing.T) {
	filesInFolder, err := os.ReadDir(testFilesDir)
	if err != nil {
		log.Fatal("i can't open the test files folder")
	}

	s, err := NewSession()
	if err != nil {
		t.Fatalf("error connecting to discord: %s", err.Error())
	}

	storage := NewStorage(s, channelId)

	for _, fileInfo := range filesInFolder {
		name := fileInfo.Name()
		file, err := os.Open(path.Join(testFilesDir, name))
		if err != nil {
			t.Fatalf("i can't open the %s test file", name)
		}

		info, err := file.Stat()
		if err != nil {
			t.Fatalf("i can't get info about the %s test file", name)
		}

		if limitDownloadSize && info.Size() > maxDownloadSize {
			t.Logf("size exceeds limit, skipping %s", name)
			continue
		}

		contents, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("i can't read the %s test file", name)
		}

		expectedSum := sha256.Sum256(contents)
		testDownload(storage, name, expectedSum, t)
	}

}

func testDownload(st DiscStorage, filename string, expectedSum [32]byte, t *testing.T) {
	buf := bytes.Buffer{}
	err := st.Receive(&buf, filename)
	if err != nil {
		t.Fatalf("error downloading '%s': %s", filename, err.Error())
	}

	checksum := sha256.Sum256(buf.Bytes())

	if expectedSum != checksum {
		t.Fatalf(
			"incorrect checksum: expected\n%s\nfound \n%s\n",
			base64.RawStdEncoding.EncodeToString(expectedSum[:]),
			base64.RawStdEncoding.EncodeToString(checksum[:]),
		)
	}
}

func NewSession() (*dg.Session, error) {
	return dg.New("Bot " + authToken)
}
