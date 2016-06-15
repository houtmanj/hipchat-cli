package cmd

import (
	"fmt"
	"net/http/httputil"
	"strings"

	"github.com/houtmanj/hipchat-cli/internal"
	"github.com/nu7hatch/gouuid"
	"github.com/spf13/cobra"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

type nagiosType struct {
	str string
}

var typeService = nagiosType{str: "service"}
var typeHost = nagiosType{str: "host"}
var typeInvalid = nagiosType{str: "invalid"}

type nagiosStatus struct {
	str   string
	color hipchat.Color
	style string
}

var statusCritical = nagiosStatus{str: "critical", color: hipchat.ColorRed, style: "lozenge-success"}
var statusWarning = nagiosStatus{str: "warning", color: hipchat.ColorYellow, style: "lozenge-current"}
var statusUnknown = nagiosStatus{str: "unknown", color: hipchat.ColorPurple, style: "lozenge-moved"}
var statusOk = nagiosStatus{str: "ok", color: hipchat.ColorGreen, style: "lozenge-success"}
var statusInvalid = nagiosStatus{str: "invalid", color: hipchat.ColorPurple, style: "lozenge-moved"}
var statusUp = nagiosStatus{str: "up", color: hipchat.ColorGreen, style: "lozenge-success"}
var statusDown = nagiosStatus{str: "down", color: hipchat.ColorRed, style: "lozenge-error"}
var statusUnreachable = nagiosStatus{str: "unreachable", color: hipchat.ColorRed, style: "lozenge-error"}

var serviceType = nagiosType{str: "service"}
var hostType = nagiosType{str: "host"}

type nagiosActions struct {
	Name string
	URL  string
}

type nagiosNotification struct {
	CheckType  nagiosType
	Status     nagiosStatus
	Service    string
	Host       string
	Output     string
	MonitorURL string
	Notify     bool
	Actions    []nagiosActions
}

// nagiosCmd represents the nagios command
var nagiosCmd = &cobra.Command{
	Use:   "nagios",
	Short: "subcommand to send nagios alerts to hipchat rooms",
	Long: `Used to send monitoring alerts to hipchat rooms.
The notifications can include linkts to various subsystems or actions.

Example:
hipchat-cli nagios --room production  --type service --status critical --service "Apache process" --output "ok - pid found" \
  --host main-web-100 --monitorurl https://nagios.com/dashboard/ \
  --actions "CreateTicket:http://jira.com"  --actions "Ack:http://nagios.com?a=ack&alert=x"
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		notif, err := validateArguments(cmd)
		if err != nil {
			return err
		}

		c, err := internal.GetClient()
		if err != nil {
			return err
		}

		n, err := getNotification(notif)
		if err != nil {
			return fmt.Errorf("error while compiling notification: %v", err)
		}

		room := cmd.Flag("room").Value.String()

		resp, err := c.Room.Notification(room, n)
		if resp != nil {
			internal.Debug(httputil.DumpResponse(resp, true))
		}
		return err
	},
}

func init() {
	RootCmd.AddCommand(nagiosCmd)

	nagiosCmd.Flags().String("type", "service", "monitoring check type: host or service")
	nagiosCmd.Flags().String("status", "", "check status: CRITICAL, WARNING, UNKNOWN, OK")

	nagiosCmd.Flags().String("service", "", "Service name")
	nagiosCmd.Flags().String("host", "", "hostname")
	nagiosCmd.Flags().String("output", "", "check output")
	nagiosCmd.Flags().String("monitorurl", "", "Url to monitoring page")
	nagiosCmd.Flags().StringSlice("actions", []string{}, "actions to put in the notification format:  <name>:<link>")
	nagiosCmd.Flags().String("room", "", "Name of the room")

	nagiosCmd.Flags().Bool("notify", false, "Send out notification to clients")
}

func validateArguments(cmd *cobra.Command) (nagiosNotification, error) {

	if cmd.Flag("room").Value.String() == "" {
		return nagiosNotification{}, fmt.Errorf("--room <room> is mandatory")
	}

	t, err := validateCheckType(cmd.Flag("type").Value.String())
	if err != nil {
		return nagiosNotification{}, err
	}

	status, err := validateStatus(cmd.Flag("status").Value.String())
	if err != nil {
		return nagiosNotification{}, err
	}

	if t == serviceType {
		if cmd.Flag("service").Value.String() == "" {
			return nagiosNotification{}, fmt.Errorf("--service is mandatory")
		}
	}

	if cmd.Flag("host").Value.String() == "" {
		return nagiosNotification{}, fmt.Errorf("--host is mandatory")
	}

	if cmd.Flag("output").Value.String() == "" {
		return nagiosNotification{}, fmt.Errorf("--output is mandatory")
	}

	actionsSlice, err := cmd.Flags().GetStringSlice("actions")
	if err != nil {
		return nagiosNotification{}, err
	}
	actions, err := validateActions(actionsSlice)
	if err != nil {
		return nagiosNotification{}, err
	}

	notif := nagiosNotification{
		CheckType:  t,
		Status:     status,
		Service:    cmd.Flag("service").Value.String(),
		Host:       cmd.Flag("host").Value.String(),
		Output:     cmd.Flag("output").Value.String(),
		MonitorURL: cmd.Flag("monitorurl").Value.String(),
		Notify:     cmd.Flag("notify").Changed,
		Actions:    actions,
	}

	return notif, nil
}

func validateCheckType(checktype string) (nagiosType, error) {
	checktype = strings.ToLower(checktype)
	switch checktype {
	case
		"service":
		return typeService, nil
	case "host":
		return typeHost, nil
	default:
		return typeInvalid, fmt.Errorf("invalid --check_type, should be service or host")
	}
}

func validateStatus(status string) (nagiosStatus, error) {
	status = strings.ToLower(status)
	switch status {
	case
		statusCritical.str:
		return statusCritical, nil
	case statusWarning.str:
		return statusWarning, nil
	case statusUnknown.str:
		return statusUnknown, nil
	case statusOk.str:
		return statusOk, nil
	case statusUp.str:
		return statusUp, nil
	case statusDown.str:
		return statusDown, nil
	case statusUnreachable.str:
		return statusUnreachable, nil
	default:
		return statusInvalid, fmt.Errorf("invalid --status, should be critical, warning, unknown or ok")
	}
}

func validateActions(actionSlice []string) (actions []nagiosActions, err error) {
	for _, v := range actionSlice {
		splits := strings.SplitN(v, ":", 2)
		if len(splits) != 2 {
			return actions, fmt.Errorf("--actions format is <action>:<URL>")
		}
		actions = append(actions, nagiosActions{splits[0], splits[1]})
	}
	return
}

func getNotification(notif nagiosNotification) (*hipchat.NotificationRequest, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	n := hipchat.NotificationRequest{
		Message: getMessage(notif),
		Notify:  notif.Notify,
		Color:   notif.Status.color,
		Card: &hipchat.Card{
			Style:       hipchat.CardStyleApplication,
			URL:         notif.MonitorURL,
			Format:      "medium",
			Title:       getTitle(notif),
			ID:          uid.String(),
			Description: hipchat.CardDescription{Format: "html", Value: notif.Output},
			Icon: &hipchat.Icon{
				URL: "https://a.fsdn.com/allura/p/nagiosplug/icon",
			},
			Attributes: getAttributeLabels(notif),
			Activity:   getActivity(notif),
		},
	}
	return &n, nil
}

func getMessage(notif nagiosNotification) string {
	if notif.CheckType == serviceType {
		return fmt.Sprintf("%v - %v on %v: %v", notif.Status.str, notif.Service, notif.Host, notif.Output)
	}
	return fmt.Sprintf("%v on %v: %v", notif.Status.str, notif.Host, notif.Output)
}

func getTitle(notif nagiosNotification) string {
	switch notif.CheckType {
	case serviceType:
		return fmt.Sprintf("Monitoring Service %v: %v on %v", notif.Status.str, notif.Service, notif.Host)
	case hostType:
		return fmt.Sprintf("Monitoring Host %v on %v", notif.Status.str, notif.Host)
	default:
		return fmt.Sprintf("Monitoring Invalid: %v on %v", notif.Status.str, notif.Host)
	}
}

func getActivity(notif nagiosNotification) *hipchat.Activity {
	switch notif.Status {
	case statusOk:
		return &hipchat.Activity{HTML: fmt.Sprintf("Recovery for %v on %v", notif.Service, notif.Host)}
	default:
		return &hipchat.Activity{HTML: fmt.Sprintf("%v for %v on %v", strings.Title(notif.Status.str), notif.Service, notif.Host)}
	}
}

func getAttributeLabels(notif nagiosNotification) []hipchat.Attribute {
	attributes := []hipchat.Attribute{}
	if notif.CheckType == serviceType {
		attributes = []hipchat.Attribute{
			getTypeLabel(notif),
			{Label: "service", Value: hipchat.AttributeValue{Label: notif.Service}},
			{Label: "host", Value: hipchat.AttributeValue{Label: notif.Host}},
		}
	} else {
		attributes = []hipchat.Attribute{
			getTypeLabel(notif),
			{Label: "host", Value: hipchat.AttributeValue{Label: notif.Host}},
		}
	}

	for _, act := range notif.Actions {
		attributes = append(attributes, hipchat.Attribute{Label: "Action", Value: hipchat.AttributeValue{Label: act.Name, URL: act.URL}})
	}

	return attributes
}

func getTypeLabel(notif nagiosNotification) hipchat.Attribute {
	return hipchat.Attribute{Label: "type", Value: hipchat.AttributeValue{Label: strings.Title(notif.Status.str), Style: notif.Status.style}}
}
