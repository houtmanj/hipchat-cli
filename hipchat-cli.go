package main

import (
	"flag"
	"fmt"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"os"
)

var secret = flag.String("secret", "", "Secret used to authenticate with hipchat")
var room = flag.String("room", "", "specify hipchat room")
var topic = flag.String("topic", "", "Change room to specified topic")
var invite = flag.String("invite", "", "User to invite to room, specify by email address")
var reason = flag.String("reason", "because..", "reason for invitation")
var msg = flag.String("msg", "", "Message to send to room")

func main() {
	flag.Parse()

	if *secret == "" {
		fmt.Println("No secret specified")
		os.Exit(1)
	}

	c := hipchat.NewClient(*secret)

	if *room == "" {
		fmt.Println("No room specified")
		os.Exit(1)
	}

	if *invite != "" {
		fmt.Printf("Inviting %v to %v with reason: %v\n", *invite, *room, *reason)
		c.Room.Invite(*room, *invite, *reason)
	} else if *msg != "" {
		fmt.Printf("Sending msg to %v: %v\n", *room, *msg)
		c.Room.Notification(*room, &hipchat.NotificationRequest{Message: *msg})
	} else if *topic != "" {
		fmt.Printf("Setting topic in %v: %v\n", *room, *topic)
		resp, err := c.Room.SetTopic(*room, *topic)
		if err != nil {
			fmt.Printf("Server returns %+v\n", resp)
			panic(err)
		}
	} else {
		fmt.Println("no action specified")
	}
}
