package main

import (
	"fmt"
	"os"
	"runtime"
)

var (
	version  = "dev"
	platform = runtime.GOOS
)

type command struct {
	name string
	run  func(args []string) error
	desc string
}

var commands = []command{
	{"install", svc, "install service init on the system"},
	{"uninstall", svc, "uninstall service init from the system"},
	{"start", svc, "start installed service"},
	{"stop", svc, "stop installed service"},
	{"restart", svc, "restart installed service"},
	{"status", svc, "return service status"},
	{"log", svc, "show service logs"},

	{"upgrade", upgrade, "upgrade the cli to the latest version"},

	{"run", run, "run the daemon"},

	{"config", cfg, "manage configuration"},

	{"activate", activation, "setup the system to use NextDNS as a resolver"},
	{"deactivate", activation, "restore the resolver configuration"},

	{"discovered", ctlCmd, "display discovered clients"},
	{"cache-stats", ctlCmd, "display cache statistics"},
	{"cache-keys", ctlCmd, "dump the list of cached entries"},
	{"trace", ctlCmd, "display a stack trace dump"},
	{"arp", ctlCmd, "dump the ARP table"},
	{"ndp", ctlCmd, "dump the NDP table"},

	{"version", showVersion, "show current version"},
}

func showCommands() {
	fmt.Println("Usage: nextdns <command> [arguments]")
	fmt.Println("")
	fmt.Println("The commands are:")
	fmt.Println("")
	for _, cmd := range commands {
		fmt.Printf("    %-15s %s\n", cmd.name, cmd.desc)
	}
	fmt.Println("")
	os.Exit(1)
}

func showVersion(args []string) error {
	fmt.Printf("nextdns version %s\n", version)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		showCommands()
	}
	cmd := os.Args[1]
	for _, c := range commands {
		if c.name != cmd {
			continue
		}
		if err := c.run(os.Args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}
	// Command not found
	showCommands()
}
