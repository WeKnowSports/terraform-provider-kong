package kong

import (
	"github.com/dghubble/sling"
)

type Config struct {
	Address string
	User string
	Password string
}

func (c *Config) Client() (*sling.Sling, error) {

	var s sling
	s = sling.new()

	if (c.User != null && c.Password != null){
		s.SetBasicAuth(c.User, c.Password)
	}

	return s.Base(c.Address), nil
}
