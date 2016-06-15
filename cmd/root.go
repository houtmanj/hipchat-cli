package cmd

import (
	"fmt"
	"os"

	"github.com/houtmanj/hipchat-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hipchat-cli",
	Short: "A CLI for common hipchat actions",
	Long: `This hipchat cli is intended to be used in scripts.

Examples of its use:
#gets topic of a room
hipchat-cli room topic --room <room>
#sets topic for a room
hipchat-cli room topic --room <room> --topic <topic>

#send notifications to a room
hipchat-cli room notify --room <room> --message <msg>

#shows a message with information regarding a monitoring alert, including handy links.
hipchat-cli nagios --room production  --type service --status critical --service "Apache process" --output "ok - pid found" \
  --host main-web-100 --monitorurl https://nagios.com/dashboard/ \
  --actions "CreateTicket:http://jira.com"  --actions "Ack:http://nagios.com?a=ack&alert=x"
`,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hipchat-cli.yaml)")
	RootCmd.PersistentFlags().BoolVar(&internal.DebugLogging, "debug", false, "Enable debugging")

	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		fmt.Println("Using specific configfile: ", cfgFile)
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".hipchat-cli") // name of config file (without extension)
	viper.AddConfigPath("$HOME")        // adding home directory as first search path
	viper.AutomaticEnv()                // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
