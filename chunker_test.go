package discordfs

import "testing"

func TestChunksNeeded(t *testing.T) {
	testCases := []struct {
		fileSize     int
		chunkSize    int
		chunksNeeded int
	}{
		{20 * KB, KB, 20},
		{20*KB + B, KB, 21},
		{512 * MB, MB, 512},
		{512*MB + 20*KB, MB, 513},
	}

	for _, testCase := range testCases {
		chunks := chunksNeeded(testCase.fileSize, testCase.chunkSize)
		if chunks != testCase.chunksNeeded {
			t.Fatalf(
				`TestChunksNeeded:
file size: %d
chunk size: %d
expected chunks needed: %d
got: %d`,
				testCase.fileSize,
				testCase.chunkSize,
				testCase.chunksNeeded,
				chunks,
			)
		}
	}
}
