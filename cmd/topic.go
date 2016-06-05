package cmd

import (
	"net/http/httputil"

	"github.com/houtmanj/hipchat-cli/internal"
	"github.com/spf13/cobra"
)

var topic string

// topicCmd represents the topic command
var topicCmd = &cobra.Command{
	Use:   "topic",
	Short: "Set or Get the topic of a room",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := internal.GetClient()
		if err != nil {
			return err
		}
		room := cmd.Flag("room").Value.String()

		if !cmd.Flag("topic").Changed {
			r, resp, err := c.Room.Get(room)

			if err != nil {
				internal.Debug(httputil.DumpResponse(resp, true))
				return err
			}

			cmd.Printf("topic: %v", r.Topic)
			return nil
		}
		cmd.Printf("Setting topic for '%v' to '%v'\n", room, topic)

		resp, err := c.Room.SetTopic(room, topic)
		if err != nil {
			internal.Debug(httputil.DumpResponse(resp, true))
			cmd.Printf("failed to set topic: %v", err)
			return err
		}

		return nil
	},
}

func init() {
	roomCmd.AddCommand(topicCmd)

	topicCmd.Flags().StringVar(&topic, "topic", "", "Specify new topic")

}
