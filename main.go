package main

import (
	zaa "./zabbix_aggregate_agent"
	"flag"
	"log"
	"strings"
)

const (
	DefaultAddress = "127.0.0.1:10050"
	DefaultTimeout = 5
)

func main() {
	var (
		listen      string
		listFile    string
		listArg     string
		listCommand string
		timeout     int
	)

	flag.StringVar(&listen, "listen", DefaultAddress, "listen address e.g. 0.0.0.0:10050")
	flag.StringVar(&listFile, "list-file", "", "zabbix-agent list file")
	flag.StringVar(&listCommand, "list-command", "", "command which prints zabbix-agent list to stdout")
	flag.StringVar(&listArg, "list", "", "zabbix-agent list , separated. e.g. 'web.example.com:10050,192.168.1.1:10050'")
	flag.IntVar(&timeout, "timeout", DefaultTimeout, "network timeout with zabbix-agent (seconds)")
	flag.Parse()

	agent := zaa.NewAgent()
	agent.Timeout = timeout

	if listFile != "" {
		agent.ListGenerator = zaa.ListFromFile
		agent.ListSource = listFile
	} else if listArg != "" {
		agent.ListGenerator = func(source string) []string {
			return strings.Split(source, ",")
		}
		agent.ListSource = listArg
	} else if listCommand != "" {
		agent.ListGenerator = zaa.ListFromCommand
		agent.ListSource = listCommand
	} else {
		log.Fatalln("option --list, --list-file or --list-command is required.")
	}

	err := agent.Run(listen)
	if err != nil {
		log.Fatalln("Error", err)
	}
}
