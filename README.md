

# hipchat-cli

hipchat-cli is a simple command line interface to common hipchat api calls.
It's main intention is to use it from scripts.


## Installation
Make sure you have a working Go environment.

To install hipchat-cli, simply run:

    > go get -v github.com/houtmanj/hipchat-cli


Make sure your `PATH` includes to the `$GOPATH/bin` directory so your commands
can be easily used:
```
export PATH=$PATH:$GOPATH/bin
```

Register the plugin with hipchat by running:
hipchat-cli register --name <nameyouwish>

This will create a small webserver and starts the flow for granting access to your hipchat.
After the flow you will be redirected to the webserver started by the command.
Here you can find the credentials of your plugin, which need to be stored in the configuration.

## configuration
~/.hipchat-cli.yaml is the default location for the configuration file.
Using --config <location> an alternative location could be specified.

Example of a configuration file:
``` yaml
---
oauthid: <secret>
oauthsecret: <secret>
proxy: http://proxy:3128
endpoint: hipchat.myspace.com
```
Only the oauthid and oauthsecret are mandatory keys.


## Examples

Playing hide and seek
```
for i in $seq 1 60); do
  hipchat-cli room notify --room hideandseek --message ${i}
  sleep 1
do

hipchat-cli room notify --room hideandseek --message "Ready or not, here I come" --notify
```

