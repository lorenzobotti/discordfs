package discordfs

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

const uploadTestSize = KB + 320*B

// delete all files that begin with "test_file"
func init() {
	st, err := newTestStorage()
	if err != nil {
		panic(err)
	}

	filesOnServer, err := st.ListFiles()
	if err != nil {
		panic(err)
	}

	for _, file := range filesOnServer {
		if strings.HasPrefix(file, "test_file") {
			st.Delete(file)
		}
	}
}

func TestUploadAndDelete(t *testing.T) {
	st, err := newTestStorage()
	if err != nil {
		panic(err)
	}

	payload := make([]byte, uploadTestSize)
	//randGen := rand.New(rand.NewSource(20021227))
	randGen := rand.New(rand.NewSource(time.Now().Unix()))
	n, err := randGen.Read(payload)
	if err != nil {
		t.Fatalf("error generating random file: %s", err.Error())
	} else if n < uploadTestSize {
		t.Fatalf("can't generate full file")
	}
	checksum := sha256.Sum256(payload)

	filename := "test_file_" + randomString(5)
	err = st.Send(bytes.NewBuffer(payload), filename, 512*B, uploadTestSize)
	if err != nil {
		t.Fatalf("error sending file: %s", err.Error())
	}

	downloader := bytes.Buffer{}
	err = st.Receive(&downloader, filename)
	if err != nil {
		t.Fatalf("error downloading file: %s", err.Error())
	}

	downloadedChecksum := sha256.Sum256(downloader.Bytes())
	if checksum != downloadedChecksum {
		t.Fatalf(
			"checksums don't match. expected\n%s\nfound\n%s\n",
			base64.RawStdEncoding.EncodeToString(checksum[:]),
			base64.RawStdEncoding.EncodeToString(downloadedChecksum[:]),
		)
	}

	err = st.Delete(filename)
	if err != nil {
		t.Fatalf("error deleting file: %s", err.Error())
	}

	filesOnServer, err := st.ListFiles()
	if err != nil {
		t.Fatalf("error getting file list: %s", err.Error())
	}

	for _, fileOnServer := range filesOnServer {
		if fileOnServer == filename {
			t.Fatalf("the file %s was evidently not deleted correctly, as it's still online", fileOnServer)
		}
	}
}

func randomString(length int) string {
	output := ""
	for i := 0; i < length; i++ {
		output += fmt.Sprintf("%c", 'a'+rand.Intn(26))
	}
	return output
}
