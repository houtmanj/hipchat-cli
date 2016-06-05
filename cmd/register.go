package cmd

import (
	"fmt"

	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httputil"
	"os/exec"
	"text/template"

	"github.com/spf13/cobra"
)

var descriptorTmpl = `{
  "description": "hipchat commandline utility",
  "key": "hipchatcli-{{.Name}}",
  "name": "{{.Name}}",
  "vendor": {
    "name": "jhoutman",
    "url": "https://github.com/houtmanj/hipchat-cli"
  },
  "links": {
    "self": "https://github.com/houtmanj/hipchat-cli/blob/master/descriptor.json"
  },
  "capabilities": {
    "hipchatApiConsumer": {
      "scopes": [
        "send_notification",
		"admin_room",
		"view_room"
      ]
    },
    "installable": {
      "allowGlobal": true,
      "allowRoom": false,
      "installedUrl": "http://localhost:8000",
      "updatedUrl": "http://localhost:8000"
    }
  }
}
`

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers the hipchat-cli as a plugin with hipchat",
	Long: `Used to register the plugin with hipchat and retrieve the required credentials.

The name of the plugin can be changed using --name. Default is: hipchat-cli
If you are using a on-premise installation then use --endpoint to point to the correct setup
`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		registerPlugin(cmd)
		startServer()
	},
}

func init() {
	RootCmd.AddCommand(registerCmd)

	registerCmd.Flags().String("name", "hipchat-cli", "Name of integration")
	registerCmd.Flags().String("endpoint", "https://www.hipchat.com", "url of hipchat server")
}

func registerPlugin(cmd *cobra.Command) {
	endpoint := cmd.Flag("endpoint").Value.String()
	name := cmd.Flag("name").Value.String()

	data := struct{ Name string }{Name: name}
	tmpl, err := template.New("test").Parse(descriptorTmpl)
	if err != nil {
		panic(err)
	}

	var doc bytes.Buffer
	err = tmpl.Execute(&doc, data)
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("%v/addons/install?url=data:application/json;base64,%v", endpoint, base64.StdEncoding.EncodeToString(doc.Bytes()))

	cmdOpen := exec.Command("open", url)
	err = cmdOpen.Run()
	if err != nil {
		fmt.Println(err)
	}

}

func printData(w http.ResponseWriter, r *http.Request) {
	r.URL.Query().Get("installable_url")
	resp, err := http.Get(r.URL.Query().Get("installable_url"))
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}

	fmt.Fprintln(w, "Store the id and secret below in the config file")

	d, _ := httputil.DumpResponse(resp, true)
	fmt.Fprintf(w, "%s", string(d))
}

func startServer() {
	http.HandleFunc("/", printData)
	http.ListenAndServe(":8000", nil)
}
