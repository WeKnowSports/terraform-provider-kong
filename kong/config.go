package kong

import (
	"github.com/dghubble/sling"
)

type Config struct {
	Address string
}

func (c *Config) Client() (*sling.Sling, error) {
	return sling.New().Base(c.Address), nil
}
