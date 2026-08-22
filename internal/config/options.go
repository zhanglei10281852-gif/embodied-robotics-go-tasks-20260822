package config

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Options struct {
	Config          Config
	ShutdownTimeout time.Duration
	AllowedOrigins  []string
}

func (c Config) Options() Options {
	return Options{Config: c, ShutdownTimeout: 10 * time.Second, AllowedOrigins: []string{"http://localhost"}}
}
func (o Options) Address() string { return net.JoinHostPort("0.0.0.0", o.Config.Port) }
func (o Options) Validate() error {
	if e := o.Config.Validate(); e != nil {
		return e
	}
	if o.ShutdownTimeout <= 0 {
		return fmt.Errorf("shutdown timeout must be positive")
	}
	for _, origin := range o.AllowedOrigins {
		if strings.TrimSpace(origin) == "" {
			return fmt.Errorf("origin cannot be empty")
		}
	}
	return nil
}
func (o Options) HasOrigin(v string) bool {
	for _, x := range o.AllowedOrigins {
		if strings.EqualFold(strings.TrimRight(x, "/"), strings.TrimRight(v, "/")) {
			return true
		}
	}
	return false
}
