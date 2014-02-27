# zabbix-aggregate-agent

zabbix-agent aggregation proxy daemon.

## Repository

[github.com/fujiwara/zabbix-aggregate-agent](https://github.com/fujiwara/zabbix-aggregate-agent)

## Build & Install

Install to $GOPATH.

    $ go get github.com/fujiwara/zabbix-aggregate-agent

## Usage

    $ zabbix-aggregate-agent --config /path/to/config.toml

## Configuration

```toml
[[agent]]
# Name: Identifier
Name = "web_servers"

# Listen: Listen address (default "127.0.0.1:10052")
Listen = "0.0.0.0:10052"

# List : List of agents' address to aggregate
List = [ "web01:10050", "web02:10050" ]

# Timeout: Timeout seconds for getting reply from agent (default 60)
Timeout = 10

# LogLevel: "Debug", "Info" or "Error" (default "Info")
LogLevel = "Debug"

[[agent]]
Name  = "app_servers"
Listen = "0.0.0.0:10053"

# ListFile: Specify the file of list of agent. ("\n" delimited)
ListFile = "/path/to/agent.list"

[[agent]]
Name = "db_servers"
Listen = "0.0.0.0:10054"

# ListCommand: Specify a command and arguments to output list of agent address. ("\n" delimited)
ListCommand = [ "/path/to/generate_list.sh", "arg1", "args2" ]

# CacheExpires : Seconds for expiring a cache of ListCommand result. (default 0 == no cache)
#   Enabled only when specify "ListCommand"
CacheExpires = 300
```

## Architecture

### Data Flow

```
[zabbix-server(or zabbix-proxy)]
    |           ^
    | (1)key    | (6)aggregated value
    |           |
    v           |
 [zabbix-aggregate-agent] <--- (2) list of zabbix-agents from static list or file or command output
    |           ^
    | (3)key    | (4)values(*)
    v           |
  [zabbix-agents(*)]
```

1. server(or proxy) requests "value" to aggregate-agent.
2. aggregate-agent resolves list of agents to forwarding request.
3. aggregate-agent requests "value" to multiple agents.
4. agnets replies "value" to aggregate-agent.
5. aggregate-agent aggregates multiple values.
  * integer or float : add
  * string : concat
6. aggregate-agent replies aggregated value to server.

## Author

Fujiwara Shunichiro <fujiwara.shunichiro@gmail.com>
