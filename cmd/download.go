package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	df "github.com/lorenzobotti/discordfs"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download one or more files from discfs",
	Long:  "Download one or more files from discfs",
	Run: func(c *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)

		if downloadCmdArgs.output != "" && len(args) > 1 {
			fmt.Fprintln(os.Stderr, "--output flag is ignored when downloading more than one file at a time")
			//	if downloadCmdArgs.force {
			//		fmt.Fprintln(os.Stderr, "continue? (y/n)")
			//
			//		scanner.Scan()
			//
			//		if !(scanner.Text() == "y" || scanner.Text() == "Y") {
			//			return
			//		}
			//	}
		}

		storage := unsafeNewStorage()
		for _, arg := range args {
			err := downloadFile(storage, scanner, arg)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	},
}

var downloadCmdArgs struct {
	output string
	force  bool
}

func downloadFile(st df.DiscStorage, scanner *bufio.Scanner, filename string) error {
	fmt.Println(filename, "starting download")
	if !downloadCmdArgs.force && fileExists(filename) {
		fmt.Println("already exists")

		fmt.Fprintf(os.Stderr, "file \"%s\" already exists. overwrite? (y/n)\n", filename)
		scanner.Scan()

		if !(scanner.Text() == "y" || scanner.Text() == "Y") {
			fmt.Println("skipping", filename)
			return nil
		}
	}

	fmt.Println("opening for writing")
	file, err := os.Create(path.Base(filename))
	if err != nil {
		return err
	}

	fmt.Println(file.Name())

	fmt.Println("starting download")
	err = st.Receive(file, filename)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, fs.ErrNotExist)
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(
		&downloadCmdArgs.output,
		"output",
		"o",
		"",
		"output file. it's ignored when downloading multiple files",
	)

	downloadCmd.Flags().BoolVarP(
		&downloadCmdArgs.force,
		"force",
		"f",
		false,
		"overwrites the existing file(s) if they already exist, otherwise it asks for each one",
	)
}
