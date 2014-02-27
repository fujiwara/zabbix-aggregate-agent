package zabbix_aggregate_agent_test

import (
	. "../zabbix_aggregate_agent"
	"testing"
	//	"reflect"
)

func TestReadConfig(t *testing.T) {
	agents, err := ReadConfig("../config_example.toml")
	if err != nil {
		t.Errorf("ReadConfig error: %v", err)
	}
	if len(agents) != 3 {
		t.Errorf("invalid agent sections")
	}
	var a0 = agents[0]
	if a0.Name != "web_servers" ||
		a0.Listen != "0.0.0.0:10052" ||
		len(a0.List) != 2 ||
		a0.List[0] != "web01:10050" ||
		a0.List[1] != "web02:10050" ||
		a0.Timeout != 10 ||
		a0.LogLevel != "Debug" {
		t.Errorf("invalid config: %v", a0)
	}
	var a1 = agents[1]
	if a1.Name != "app_servers" ||
		a1.Listen != "0.0.0.0:10053" ||
		a1.ListFile != "/path/to/agent.list" {
		t.Errorf("invalid config: %v", a1)
	}
	var a2 = agents[2]
	if a2.Name != "db_servers" ||
		a2.Listen != "0.0.0.0:10054" ||
		len(a2.ListCommand) != 3 ||
		a2.ListCommand[0] != "/path/to/generate_list.sh" ||
		a2.ListCommand[1] != "arg1" ||
		a2.ListCommand[2] != "args2" {
		t.Errorf("invalid config: %v", a2)
	}
}

func TestReadConfigFail(t *testing.T) {
	configAgents, err := ReadConfig("config_example.tom")
	if err == nil {
		t.Errorf("no errors in ReadConfig")
	}
	if len(configAgents) > 0 {
		t.Errorf("must be returned no configAgents")
	}
}
