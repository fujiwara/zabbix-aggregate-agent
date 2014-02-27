package zabbix_aggregate_agent

import (
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

type Cache struct {
	List      []string
	UpdatedAt time.Time
	Expires   time.Duration
	Generator func() ([]string, error)
}

func NewCachedListGenerator(generator func() ([]string, error), expires int) (cachedGenerator func() ([]string, error)) {
	if expires <= 0 {
		cachedGenerator = generator
		return
	}
	cache := &Cache{
		Expires:   time.Duration(int64(expires)) * time.Second,
		Generator: generator,
	}
	cachedGenerator = func() ([]string, error) {
		var err error
		now := time.Now()
		expired := cache.UpdatedAt.Add(cache.Expires)
		if now.After(expired) {
			cache.List, err = cache.Generator()
			cache.UpdatedAt = now
		}
		return cache.List, err
	}
	return cachedGenerator
}

func NewListFromArgGenerator(source string) (f func() ([]string, error)) {
	f = func() (list []string, err error) {
		list = listFromString(source, ",")
		return
	}
	return
}

func NewListFromFileGenerator(filename string) (f func() ([]string, error)) {
	f = func() (list []string, err error) {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return
		}
		contentString := string(content)
		list = listFromString(contentString, "\n")
		return
	}
	return
}

func NewListFromCommandGenerator(command string) (f func() ([]string, error)) {
	f = func() (list []string, err error) {
		out, err := exec.Command(command).Output()
		if err != nil {
			return
		}
		outString := string(out)
		list = listFromString(outString, "\n")
		return
	}
	return
}

func listFromString(content string, delimiter string) (list []string) {
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
