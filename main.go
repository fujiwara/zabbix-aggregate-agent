package main

import (
	zaa "./zabbix_aggregate_agent"
	"flag"
	"log"
)

func runByConfig(config string) {
	agents, err := zaa.NewAgentsFromConfig(config)
	if err != nil {
		log.Fatalln(err)
	}
	ch := make(chan bool)
	for _, agent := range agents {
		go agent.RunNotify(ch)
	}
	for _, _ = range agents {
		<-ch
	}
	log.Fatalln("All of agents could not be run.")
}

func main() {
	var (
		listen      string
		listFile    string
		listArg     string
		listCommand string
		timeout     int
		expires     int
		config      string
	)
	flag.StringVar(&listen, "listen", zaa.DefaultAddress, "listen address e.g. 0.0.0.0:10052")
	flag.StringVar(&listFile, "list-file", "", "zabbix-agent list file")
	flag.StringVar(&listCommand, "list-command", "", "command which prints zabbix-agent list to stdout")
	flag.StringVar(&listArg, "list", "", "zabbix-agent list , separated. e.g. 'web.example.com:10050,192.168.1.1:10050'")
	flag.IntVar(&timeout, "timeout", 0, "network timeout with zabbix-agent (seconds)")
	flag.IntVar(&expires, "expires", 0, "list cache expires (seconds)")
	flag.StringVar(&config, "config", "", "config file")
	flag.Parse()

	if config != "" {
		runByConfig(config)
		return
	}

	agent := zaa.NewAgent("1", listen, timeout)
	if listFile != "" {
		agent.ListGenerator = zaa.NewListFromFileGenerator(listFile)
	} else if listArg != "" {
		agent.ListGenerator = zaa.NewListFromArgGenerator(listArg)
	} else if listCommand != "" {
		agent.ListGenerator = zaa.NewCachedListGenerator(
			zaa.NewListFromCommandGenerator(listCommand),
			expires,
		)
	} else {
		log.Fatalln("option either --list, --list-file or --list-command is required.")
	}
	err := agent.Run()
	if err != nil {
		log.Fatalln("Error", err)
	}
	return
}
