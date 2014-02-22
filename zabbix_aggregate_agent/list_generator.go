package zabbix_aggregate_agent

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func ListFromFile(filename string) (list []string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return
	}
	list = ListFromString(string(content))
	return
}

func ListFromCommand(command string) (list []string) {
	out, err := exec.Command(command).Output()
	if err != nil {
		log.Println(err)
		return
	}
	list = ListFromString(string(out))
	return
}

func ListFromString(content string) (list []string) {
	n := 0
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Index(line, "#") == 0 || line == "" {
			// comment out
			continue
		}
		n++
		list = append(list, line)
	}
	return
}
