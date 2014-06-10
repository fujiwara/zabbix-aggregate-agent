package main

import (
	zaa "github.com/fujiwara/zabbix-aggregate-agent/zabbix_aggregate_agent"
	"flag"
	"log"
)

func main() {
	var config string
	flag.StringVar(&config, "config", "", "config file")
	flag.Parse()
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
