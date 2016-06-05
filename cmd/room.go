package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// roomCmd represents the room command
var roomCmd = &cobra.Command{
	Use:   "room",
	Short: "Perform actions on a hipchat room",
	Long: `Allows you to perform common operations on a room. For example:

notify: Send a message to a room
topic:  get or set the topic
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flag("room").Changed {
			return fmt.Errorf("Specification of a room is mandatory, use --name")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(roomCmd)

	roomCmd.PersistentFlags().String("room", "", "Name of the room")
}
