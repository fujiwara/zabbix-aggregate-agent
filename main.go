package main

import (
	zaa "github.com/fujiwara/zabbix-aggregate-agent/zabbix_aggregate_agent"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {
	var config string
	var showVersion bool

	flag.StringVar(&config, "config", "", "config file")
	flag.BoolVar(&showVersion, "version", false, "show Version")
	flag.Parse()

	if showVersion {
		fmt.Printf("zabbix-aggregate-agent version %s (revision %s)\n", Version, Revision)
		os.Exit(255)
	}

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
