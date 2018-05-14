package kong

import (
	"github.com/dghubble/sling"
)

type Config struct {
	Address  string
	Username string
	Password string
	JWT string
}

func (c *Config) Client() (*sling.Sling, error) {
	if c.JWT != "" {
		return sling.New().Set("Authorization", "Bearer " + c.JWT).Base(c.Address), nil
	} else {
        return sling.New().SetBasicAuth(c.Username, c.Password).Base(c.Address), nil
	}
}
