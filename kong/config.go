package kong

import (
	"github.com/dghubble/sling"
)

type Config struct {
	Address string
	Username string
	Password string
	Headers map[string]interface{}
}

func (c *Config) Client() (*sling.Sling, error) {
	client := sling.New().SetBasicAuth(c.Username, c.Password).Base(c.Address)

	for key, value := range c.Headers {
    client.Set(key, value.(string))
	}

	return client, nil
}
