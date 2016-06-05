package internal

import (
	"fmt"

	"net/http/httputil"

	"github.com/spf13/viper"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

//GetClient return a hipchat client
//it used oauthid and oauthsecret to retrieve a temporary access token
//It also listens to proxy and enpoint configuration in the configfile.
func GetClient() (*hipchat.Client, error) {

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

	proxy := viper.GetString("proxy")
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, fmt.Errorf("Could not determine proxy URL: %v", err)
		}
		proxyClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
		c.SetHTTPClient(proxyClient)
	}

	endpoint := viper.GetString("endpoint")
	if endpoint != "" {
		endpointURL, err := url.Parse(endpoint)
		if err != nil {
			return nil, fmt.Errorf("Could not determine endpoint URL: %v", err)
		}
		c.BaseURL = endpointURL
	}

	token, resp, err := c.GenerateToken(hipchat.ClientCredentials{oauthID, oauthSecret}, scope)
	Debug(httputil.DumpResponse(resp, true))

	if err != nil {
		return nil, fmt.Errorf("Error while retrieving oath token: %v", err)
	}

	c = token.CreateClient()
	return c, nil
}
