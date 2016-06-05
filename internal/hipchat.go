package internal

import (
	"fmt"

	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/spf13/viper"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

//GetClient return a hipchat client
//it used oauthid and oauthsecret to retrieve a temporary access token
//It also listens to proxy and enpoint configuration in the configfile.
func GetClient() (*hipchat.Client, error) {

	//changing the default http client, because the library does not allways
	//use the passed in client
	err := configDefaultHTTPClient()
	if err != nil {
		return nil, fmt.Errorf("Error while configuring default httpclient: %v", err)
	}

	oauthID := viper.GetString("oauthid")
	if oauthID == "" {
		return nil, fmt.Errorf("Specify an oauthid in the config file")
	}

	oauthSecret := viper.GetString("oauthsecret")
	if oauthSecret == "" {
		return nil, fmt.Errorf("Specify an oauthsecret in the config file")
	}

	scope := []string{"admin_room", "view_room", "send_notification"}
	c := hipchat.NewClient("")
	c, err = configClient(c)
	if err != nil {
		return nil, fmt.Errorf("Error while configuring client: %v", err)
	}

	token, resp, err := c.GenerateToken(hipchat.ClientCredentials{oauthID, oauthSecret}, scope)
	if resp != nil {
		Debug(httputil.DumpResponse(resp, true))
	}
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving oath token: %v", err)
	}

	c = token.CreateClient()
	c, err = configClient(c)
	if err != nil {
		return nil, fmt.Errorf("Error while configuring client: %v", err)
	}
	return c, nil
}

func configDefaultHTTPClient() error {
	proxy := viper.GetString("proxy")
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return fmt.Errorf("Could not determine proxy URL: %v", err)
		}

		if DebugLogging {
			fmt.Println("Using proxy: ", proxyURL.String())
		}

		proxyClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
		http.DefaultClient = proxyClient
	}
	if DebugLogging {
		fmt.Println("Not using a proxy.")
	}
	return nil
}

func configClient(c *hipchat.Client) (*hipchat.Client, error) {
	endpoint := viper.GetString("endpoint")
	if endpoint != "" {
		endpointURL, err := url.Parse(endpoint)
		if err != nil {
			return nil, fmt.Errorf("Could not determine endpoint URL: %v", err)
		}
		c.BaseURL = endpointURL
	}

	return c, nil
}
