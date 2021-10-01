package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "discfs",
	Short: "Use a Discord channel as your own private, unlimited cloud storage",
	Long:  "Use a Discord channel as your own private, unlimited cloud storage",
	//Run: func(c *cobra.Command, args []string) {
	//	for _, arg := range args {
	//		fmt.Println(arg)
	//	}
	//},
}

func main() {
	rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.main.yaml")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".disc-fs")
	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}

func getTokenAndChannel() (string, string, error) {
	token, ok := viper.Get("token").(string)
	if !ok {
		return "", "", errors.New("token wasn't provided")
	}
	channel, ok := viper.Get("channel").(string)
	if !ok {
		return "", "", errors.New("token wasn't provided")
	}

	return token, channel, nil
}
