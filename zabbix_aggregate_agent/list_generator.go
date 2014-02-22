package zabbix_aggregate_agent

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func ListFromArg (source string) (list []string) {
	list = ListFromString(source, ",")
	return
}

func ListFromFile(filename string) (list []string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return
	}
	list = ListFromString(string(content), "\n")
	return
}

func ListFromCommand(command string) (list []string) {
	out, err := exec.Command(command).Output()
	if err != nil {
		log.Println(err)
		return
	}
	list = ListFromString(string(out), "\n")
	return
}

func ListFromString(content string, delimiter string) (list []string) {
	n := 0
	lines := strings.Split(content, delimiter)
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
