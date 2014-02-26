package zabbix_aggregate_agent

import (
	"fmt"
	"io/ioutil"
	"log"
	"github.com/BurntSushi/toml"
)

type agent struct {
	Name        string
	Listen      string
	ListFile    string
	ListCommand string
	List        string
	Timeout     int
	Expires     int
}

type agents struct {
	Agent []agent
}

func BuildAgentsFromConfig (filename string) (agentInstances []*Agent, err error) {
	log.Println("Loading config file:", filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	log.Println("\n", string(content))
	var config agents
	if _, err = toml.Decode(string(content), &config); err != nil {
		return
	}
	for i, c := range config.Agent {
		if c.Name == "" {
			c.Name = fmt.Sprintf("%d", i + 1)
		}
		log.Println("Defining agent", c.Name)
		instance := NewAgent(c.Name, c.Listen, c.Timeout)
		if c.ListFile != "" {
			instance.ListGenerator = ListFromFile
			instance.ListSource    = c.ListFile
		} else if c.List != "" {
			instance.ListGenerator = ListFromArg
			instance.ListSource    = c.List
		} else if c.ListCommand != "" {
			instance.ListGenerator = CachedListGenerator(ListFromCommand, c.Expires)
			instance.ListSource = c.ListCommand
		} else {
			log.Fatalln("option List, ListFile or ListCommand is required.")
		}
		agentInstances = append(agentInstances, instance)
	}
	return
}
