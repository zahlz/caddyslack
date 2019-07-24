package caddyslack

import (
	"errors"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("slack", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

//Config describes all options that can be set for the plugin.
type config struct {
	endpoint  string
	remoteURL string
	delete    []string
	only      []string
}

func newConfig() *config {
	return &config{endpoint: "/slack"}
}

func setup(c *caddy.Controller) error {
	sc, err := parse(c)
	if err != nil {
		return err
	}

	if c.ServerBlockKeyIndex == 0 {
		c.ServerBlockStorage = newHandler(sc)
	}

	if slackHandler, ok := c.ServerBlockStorage.(*handler); ok {
		httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
			slackHandler.Next = next
			return slackHandler
		})
		return nil
	}
	return errors.New("[slack] Could not create the middleware handler")
}

func parse(c *caddy.Controller) (conf *config, err error) {
	conf = newConfig()
	for c.Next() {
		args := c.RemainingArgs()

		switch len(args) {
		case 1:
			conf.endpoint = args[0]
		}

		err := iterateBlocks(c, conf)
		if err != nil {
			return nil, err
		}
	}

	if conf.remoteURL == "" {
		return nil, errors.New("required field 'url' not found")
	}
	return conf, nil
}

func iterateBlocks(c *caddy.Controller, conf *config) error {
	for c.NextBlock() {
		switch c.Val() {
		case "url":
			if !c.NextArg() {
				return c.ArgErr()
			}
			conf.remoteURL = c.Val()
		case "delete":
			for c.NextBlock() {
				conf.delete = append(conf.delete, c.Val())
			}
		case "only":
			conf.only = make([]string, 0)
			for c.NextBlock() {
				conf.only = append(conf.only, c.Val())
			}
		}
	}
	return nil
}
