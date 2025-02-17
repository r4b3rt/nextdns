package ubios

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/nextdns/nextdns/config"
	"github.com/nextdns/nextdns/router/internal"
)

type Router struct {
	DNSMasqPath     string
	ListenPort      string
	ClientReporting bool
}

func isUnifi() bool {
	if st, _ := os.Stat("/data/unifi"); st != nil && st.IsDir() {
		return true
	}
	if err := exec.Command("ubnt-device-info", "firmware").Run(); err == nil {
		return true
	}
	return false
}

func New() (*Router, bool) {
	if !isUnifi() {
		return nil, false
	}
	return &Router{
		DNSMasqPath: "/run/dnsmasq.conf.d/nextdns.conf",
		ListenPort:  "5342",
	}, true
}

func (r *Router) String() string {
	return "ubios"
}

func (r *Router) Configure(c *config.Config) error {
	if dnsFilterEnabled() {
		return fmt.Errorf(`UDM "Content Filtering" feature is enabled. Please disable it to use NextDNS`)
	}
	c.Listens = []string{net.JoinHostPort("localhost", r.ListenPort)}
	r.ClientReporting = c.ReportClientInfo
	if c.CacheSize == "0" || c.CacheSize == "" {
		// Make sure we setup a non-0 cache as we disable dnsmasq cache
		c.CacheSize = "10MB"
	}
	return nil
}

func (r *Router) Setup() error {
	return r.setupDNSMasq()
}

func (r *Router) Restore() error {
	if err := os.Remove(r.DNSMasqPath); err != nil {
		return err
	}
	return killDNSMasq()
}

func (r *Router) setupDNSMasq() error {
	if err := internal.WriteTemplate(r.DNSMasqPath, tmpl, r, 0644); err != nil {
		return err
	}
	return killDNSMasq()
}

func dnsFilterEnabled() bool {
	_, err := os.Stat("/run/dnsfilter/dnsfilter")
	return err == nil
}

func killDNSMasq() error {
	b, err := os.ReadFile("/run/dnsmasq.pid")
	if err != nil {
		return err
	}
	pid := string(bytes.TrimSpace(b))
	if err := exec.Command("kill", pid).Run(); err != nil {
		return fmt.Errorf("dnsmasq kill: %v", err)
	}
	return nil
}

var tmpl = `# Configuration generated by NextDNS
no-resolv
server=127.0.0.1#{{.ListenPort}}
{{- if .ClientReporting}}
add-mac
{{- end}}
add-subnet=32,128
max-cache-ttl=0
`
