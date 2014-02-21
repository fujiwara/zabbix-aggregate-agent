package zabbix_aggregate_agent

import (
	"os/exec"
	"io/ioutil"
	"strings"
	"log"
)

func ListFromFile(filename string) (list []string, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return
	}
	list, err = ListFromString(string(content))
	return
}

func ListFromCommand(command string) (list []string, err error) {
	out, err := exec.Command(command).Output()
	if err != nil {
		log.Println(err)
		return
	}
	list, err = ListFromString(string(out))
	return
}

func ListFromString(content string) (list []string, err error) {
	n := 0
	lines := strings.Split(content, "\n")
	for _, line := range(lines) {
		if strings.Index(line, "#") == 0 || line == "" {
			// comment out
			continue
		}
		n++
		list = append(list, line)
	}
	return
}
