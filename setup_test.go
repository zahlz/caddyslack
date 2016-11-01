package caddyslack

import (
	"errors"
	"testing"

	"github.com/mholt/caddy"
	"github.com/stretchr/testify/assert"
)

func TestSetupParse(t *testing.T) {
	tests := []struct {
		config     string
		expectErr  error
		expectConf func() *config
	}{
		{
			`slack`,
			errors.New("required field 'url' not found"),
			nil,
		},
		{`slack /slacking`,
			errors.New("required field 'url' not found"),
			nil,
		},
		{
			`slack {
		   url https://hooks.slack.com/services/ID/TOKEN
		  }`,
			nil,
			func() *config {
				return &config{endpoint: "/slack", remoteURL: "https://hooks.slack.com/services/ID/TOKEN"}
			},
		},
		{
			`slack /slacking {
		   url https://hooks.slack.com/services/ID/TOKEN
		  }`,
			nil,
			func() *config {
				return &config{endpoint: "/slacking", remoteURL: "https://hooks.slack.com/services/ID/TOKEN"}
			},
		},
		{
			`slack {
				url
			}`,
			errors.New("Testfile:2 - Parse error: Wrong argument count or unexpected line ending after 'url'"),
			nil,
		},
		{
			`slack {
		   url https://hooks.slack.com/services/ID/TOKEN
			 delete
			 	 channel
		  }`,
			nil,
			func() *config {
				return &config{endpoint: "/slack", remoteURL: "https://hooks.slack.com/services/ID/TOKEN", delete: []string{"channel"}}
			},
		},
		{
			`slack {
		   url https://hooks.slack.com/services/ID/TOKEN
			 delete
			 	 abc.xyz.channel
		  }`,
			nil,
			func() *config {
				return &config{endpoint: "/slack", remoteURL: "https://hooks.slack.com/services/ID/TOKEN", delete: []string{"abc.xyz.channel"}}
			},
		},
		{
			`slack {
		   url https://hooks.slack.com/services/ID/TOKEN
			 delete
			 	 channel
				 text
		  }`,
			nil,
			func() *config {
				return &config{endpoint: "/slack", remoteURL: "https://hooks.slack.com/services/ID/TOKEN", delete: []string{"channel", "text"}}
			},
		},
	}

	for i, test := range tests {
		c := caddy.NewTestController("http", test.config)
		mc, err := parse(c)
		if test.expectErr != nil {
			assert.Nil(t, mc, "Index %d", i)
			assert.EqualError(t, err, test.expectErr.Error(), "Index %d with config:\n%s", i, test.config)
			continue
		}
		assert.NoError(t, err, "Index %d", i)
		assert.Exactly(t, test.expectConf(), mc, "Index %d", i)
	}
}
