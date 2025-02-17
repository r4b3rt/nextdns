package generic

import (
	"github.com/nextdns/nextdns/config"
)

type Router struct {
}

func New() *Router {
	return &Router{}
}

func (r *Router) String() string {
	return "generic"
}

func (r *Router) Configure(c *config.Config) error {
	c.Listens = []string{":53"}
	return nil
}

func (r *Router) Setup() error {
	return nil
}

func (r *Router) Restore() error {
	return nil
}
