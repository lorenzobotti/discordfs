package main

import (
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	df "github.com/lorenzobotti/discordfs"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload one or more files to discfs",
	Long:  "Upload one or more files to discfs",
	Run: func(c *cobra.Command, args []string) {
		for _, arg := range args {
			uploadFile(arg)
		}
	},
}

func uploadFile(filename string) error {
	storage := unsafeNewStorage()

	file, err := os.Open(filename)
	cobra.CheckErr(err)

	stat, err := file.Stat()
	cobra.CheckErr(err)

	err = storage.Send(file,
		df.NewFileInfo(
			path.Join(uploadCmdArgs.path, path.Base(filename)),
			stat.ModTime(), int(stat.Size()),
		),
		getChunkSize(),
	)
	cobra.CheckErr(err)

	return nil
}

var uploadCmdArgs struct {
	chunkSize string
	path      string
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVar(&uploadCmdArgs.chunkSize,
		"chunk",
		"4MB",
		`size of each chunk the file gets split in. max allowed by Discord is 8MB.
sizes larger than 4MB may fail, we're working on it`,
	)

	uploadCmd.Flags().StringVar(&uploadCmdArgs.path,
		"path",
		"/",
		"remote upload path (only the folder name)",
	)
}

func getChunkSize() int {
	size, err := parseFileSize(uploadCmdArgs.chunkSize)
	cobra.CheckErr(err)

	return size
}

func parseFileSize(in string) (int, error) {
	names := map[string]int{
		//"B":  df.B,
		"KB": df.KB,
		"MB": df.MB,
		"GB": df.GB,
		"TB": df.TB,
	}

	for name, nameSize := range names {
		if !strings.HasSuffix(in, name) {
			continue
		}

		sizeArg := strings.TrimSuffix(in, name)
		size, err := strconv.Atoi(sizeArg)

		return nameSize * size, err
	}

	return 0, errors.New("can't parse file size")
}
