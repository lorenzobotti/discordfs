package main

import (
	"fmt"
	"time"

	df "github.com/lorenzobotti/discordfs"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list [path]",
	Aliases: []string{"ls"},
	Short:   "List the files on the server",
	Long:    "List the files and folders on the server at the specified path",
	Run: func(c *cobra.Command, args []string) {
		listCmdArgs.folder = df.Root
		if len(args) >= 1 {
			listCmdArgs.folder = args[0]
		}

		listFiles()
	},
}

var listCmdArgs struct {
	color    bool
	metadata bool
	folder   string
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(
		&listCmdArgs.color,
		"color",
		"c",
		true,
		"color the folders' names",
	)
	listCmd.Flags().BoolVarP(&listCmdArgs.metadata,
		"all",
		"a",
		false,
		"print file metadata",
	)
}

func listFiles() {
	st := unsafeNewStorage()

	files, err := st.ListFiles(listCmdArgs.folder)
	cobra.CheckErr(err)

	for _, file := range files {
		if listCmdArgs.metadata {
			if file.IsDir() {
				fmt.Print("d ")
			} else {
				fmt.Print("  ")
			}

			fmt.Print(file.ModTime().Format(time.Stamp) + " ")
		}

		if file.IsDir() && listCmdArgs.color {
			fmt.Println(ansiColorBlue(file.Name()))
		} else {
			fmt.Println(file.Name())
		}
	}
}

func ansiColorBlue(in string) string {
	return "\x1b[34m" + in + "\x1b[0m"
}
