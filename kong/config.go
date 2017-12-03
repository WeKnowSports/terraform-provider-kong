package kong

import (
	"github.com/dghubble/sling"
)

type Config struct {
	Address  string
	Username string
	Password string
}

func (c *Config) Client() (*sling.Sling, error) {
	return sling.New().SetBasicAuth(c.Username, c.Password).Base(c.Address), nil
}
