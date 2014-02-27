package zabbix_aggregate_agent

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"strings"
)

type agent struct {
	Name         string
	Listen       string
	ListFile     string
	ListCommand  []string
	List         []string
	Timeout      int
	CacheExpires int
	LogLevel     string
}

type agents struct {
	Agent []agent
}

func BuildAgentsFromConfig(filename string) (agentInstances []*Agent, err error) {
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
			c.Name = fmt.Sprintf("%d", i+1)
		}
		log.Println("Defining agent", c.Name)
		instance := NewAgent(c.Name, c.Listen, c.Timeout)
		if c.ListFile != "" {
			instance.ListGenerator = NewListFromFileGenerator(c.ListFile)
		} else if len(c.List) > 0 {
			instance.ListGenerator = NewListGenerator(c.List)
		} else if len(c.ListCommand) > 0 {
			command := c.ListCommand[0]
			args := c.ListCommand[1:]
			instance.ListGenerator = NewCachedListGenerator(
				NewListFromCommandGenerator(command, args...),
				c.CacheExpires,
			)
		} else {
			log.Fatalln("option List, ListFile or ListCommand is required.")
		}

		switch strings.ToUpper(c.LogLevel) {
		case "DEBUG":
			instance.MinLogLevel = Debug
		case "INFO":
			instance.MinLogLevel = Info
		case "ERROR":
			instance.MinLogLevel = Error
		case "":
			// default
		default:
			log.Println("LogLevel", c.LogLevel, "is unsupported. Using default level", LogLabel[DefaultLogLevel])
		}

		agentInstances = append(agentInstances, instance)
	}
	return
}
