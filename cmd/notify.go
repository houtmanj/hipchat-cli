package cmd

import (
	"fmt"

	"net/http/httputil"

	"github.com/houtmanj/hipchat-cli/internal"
	"github.com/tbruyelle/hipchat-go/hipchat"

	"github.com/spf13/cobra"
)

var notify bool

// notifyCmd represents the notify command
var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Sends a notification to a room",
	Long:  `Use the --notify to indicate if the room members should receive a notification`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flag("message").Changed {
			return fmt.Errorf("no message specified, use --message")
		}
		room := cmd.Flag("room").Value.String()
		message := cmd.Flag("message").Value.String()

		c, err := internal.GetClient()
		if err != nil {
			return err
		}

		cmd.Printf("Sending '%v' to %v\n", message, room)

		resp, err := c.Room.Notification(room, &hipchat.NotificationRequest{Message: message, Notify: notify})
		internal.Debug(httputil.DumpResponse(resp, true))
		return err
	},
}

func init() {
	roomCmd.AddCommand(notifyCmd)

	notifyCmd.Flags().String("message", "", "Message to send")
	notifyCmd.Flags().BoolVar(&notify, "notify", false, "Send out notification to clients")

}
