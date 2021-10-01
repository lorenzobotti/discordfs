package main

import (
	dg "github.com/bwmarrin/discordgo"
	df "github.com/lorenzobotti/discordfs"
	"github.com/spf13/cobra"
)

func unsafeNewStorage() df.DiscStorage {
	token, channel, err := getTokenAndChannel()
	cobra.CheckErr(err)

	session, err := dg.New("Bot " + token)
	cobra.CheckErr(err)

	return df.NewStorage(session, channel)
}
