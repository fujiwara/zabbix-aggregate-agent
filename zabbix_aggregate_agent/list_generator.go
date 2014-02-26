package zabbix_aggregate_agent

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

type Cache struct {
	List      []string
	UpdatedAt time.Time
	Expires   time.Duration
	Generator func(string) []string
}

func CachedListGenerator(generator func(string) []string, expires int) (cachedGenerator func(string) []string) {
	if expires <= 0 {
		return generator
	}
	cache := &Cache{
		Expires:   time.Duration(int64(expires)) * time.Second,
		Generator: generator,
	}
	cachedGenerator = func(source string) []string {
		now := time.Now()
		expired := cache.UpdatedAt.Add(cache.Expires)
		if now.After(expired) {
			cache.List = cache.Generator(source)
			cache.UpdatedAt = now
		}
		return cache.List
	}
	return cachedGenerator
}

func ListFromArg(source string) (list []string) {
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
	log.Println("invoking command:", command)
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
