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

func TestUploadAndDelete(t *testing.T) {
	st, err := newTestStorage()
	if err != nil {
		panic(err)
	}

	err = deleteTestFiles(st)
	if err != nil {
		t.Fatalf("error calling deleteTestFiles(): %s", err.Error())
	}

	var uploadTestSizes = []int{
		KB + 320*B,
		KB,
		B,
		36 * B,
	}

	chunkSize := 512 * B

	for _, size := range uploadTestSizes {
		randomFileUploadChecksumRemove(st, size, chunkSize, t)
	}
}

// full testing of the upload process:
// generate a random data buffer,
// download it back,
// check it's correct,
// delete it,
// check it's not there anymore.
func randomFileUploadChecksumRemove(st DiscStorage, uploadTestSize, chunkSize int, t *testing.T) {
	payload := make([]byte, uploadTestSize)

	// fill with random bytes
	randGen := rand.New(rand.NewSource(time.Now().Unix()))
	n, err := randGen.Read(payload)
	if err != nil {
		t.Fatalf("error generating random file: %s", err.Error())
	} else if n < uploadTestSize {
		t.Fatalf("can't generate full file")
	}
	checksum := sha256.Sum256(payload)

	filename := "test_file_" + randomString(5) + ".txt"
	err = st.Send(bytes.NewBuffer(payload), filename, chunkSize, uploadTestSize)
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

	fileExists, err := st.DoesFileExist(filename)
	if err != nil {
		t.Fatalf("error checking the file's presence: %s", err.Error())
	}

	if !fileExists {
		t.Fatal("the file was uploaded and downloaded correctly, but st.DoesFileExist() can't find it")
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
		if fileOnServer.Name() == filename {
			t.Fatalf("the file %s was evidently not deleted correctly, as it's still online", fileOnServer.Name())
		}
	}
}

// delete all files that begin with "test_file"
func deleteTestFiles(st DiscStorage) error {
	filesOnServer, err := st.ListFiles()
	if err != nil {
		return err
	}

	for _, file := range filesOnServer {
		if strings.HasPrefix(file.Name(), "test_file") {
			err = st.Delete(file.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func randomString(length int) string {
	output := ""
	for i := 0; i < length; i++ {
		output += fmt.Sprintf("%c", 'a'+rand.Intn(26))
	}
	return output
}
