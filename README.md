# zabbix-aggregate-agent

zabbix-agent aggregation proxy daemon.

## Repository

[github.com/fujiwara/zabbix-aggregate-agent](https://github.com/fujiwara/zabbix-aggregate-agent)

## Getting binary

 * [Mac OS X i386](http://fujiwara.github.io/bin/darwin-386/zabbix-aggregate-agent)
 * [Mac OS X x86_64](http://fujiwara.github.io/bin/darwin-amd64/zabbix-aggregate-agent)
 * [Linux i386](http://fujiwara.github.io/bin/linux-386/zabbix-aggregate-agent)
 * [Linux x86_64](http://fujiwara.github.io/bin/linux-amd64/zabbix-aggregate-agent)
 * [Linux arm](http://fujiwara.github.io/bin/linux-arm/zabbix-aggregate-agent)
 * [Windows i386](http://fujiwara.github.io/bin/windows-386/zabbix-aggregate-agent)
 * [Windows x86_64](http://fujiwara.github.io/bin/windows-amd64/zabbix-aggregate-agent)

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


## Example

### start

```
$ zabbix-aggregate-agent --config config.toml
2014/02/27 19:48:19 Loading config file: config.toml
----------------------------------
[[agent]]
Name     = "example"
List     = ["10.8.0.1:10050", "10.8.0.2:10050"]
LogLevel = "Debug"

----------------------------------
2014/02/27 19:48:19 Initialize agent example
2014/02/27 19:48:19 INFO [example] Listing 127.0.0.1:10052
2014/02/27 19:48:19 INFO [example] Ready for connection
```

### get "system.uptime" (integer)

```
$ zabbix_get -s 127.0.0.1 -p 10052 -k "system.uptime"
66407369

2014/02/27 19:48:49 DEBUG [example] Accepted connection from 127.0.0.1:65038
2014/02/27 19:48:49 DEBUG [example] Key: system.uptime
2014/02/27 19:48:49 DEBUG [example] List: [10.8.0.1:10050 10.8.0.2:10050]
2014/02/27 19:48:49 DEBUG [example] Sending key: system.uptime to 10.8.0.1:10050
2014/02/27 19:48:49 DEBUG [example] Sending key: system.uptime to 10.8.0.2:10050
2014/02/27 19:48:49 DEBUG [example] Replied from 10.8.0.1:10050 in 25 msec: 53034117
2014/02/27 19:48:49 DEBUG [example] Replied from 10.8.0.2:10050 in 71 msec: 13373252
2014/02/27 19:48:49 DEBUG [example] Aggregated system.uptime = 66407369
2014/02/27 19:48:49 DEBUG [example] Closing connection: 127.0.0.1:65038
```

### get "system.cpu.load[]" (float)

```
$ zabbix_get -s 127.0.0.1 -p 10052 -k "system.cpu.load[]"
0.360000

2014/02/27 19:55:38 DEBUG [example] Accepted connection from 127.0.0.1:65080
2014/02/27 19:55:38 DEBUG [example] Key: system.cpu.load[]
2014/02/27 19:55:38 DEBUG [example] List: [10.8.0.1:10050 10.8.0.2:10050]
2014/02/27 19:55:38 DEBUG [example] Sending key: system.cpu.load[] to 10.8.0.1:10050
2014/02/27 19:55:38 DEBUG [example] Sending key: system.cpu.load[] to 10.8.0.2:10050
2014/02/27 19:55:38 DEBUG [example] Replied from 10.8.0.1:10050 in 26 msec: 0.160000
2014/02/27 19:55:38 DEBUG [example] Replied from 10.8.0.2:10050 in 71 msec: 0.200000
2014/02/27 19:55:38 DEBUG [example] Aggregated system.cpu.load[] = 0.360000
2014/02/27 19:55:38 DEBUG [example] Closing connection: 127.0.0.1:65080
```

### get "system.uname" (string)

```
$ zabbix_get -s 127.0.0.1 -p 10052 -k "system.uname"
Linux www 3.2.0-25-generic #40-Ubuntu SMP Wed May 23 20:30:51 UTC 2012 x86_64 x86_64 x86_64 GNU/Linux
Linux raspberrypi 3.6.11+ #474 PREEMPT Thu Jun 13 17:14:42 BST 2013 armv6l GNU/Linux

2014/02/27 19:48:55 DEBUG [example] Accepted connection from 127.0.0.1:65041
2014/02/27 19:48:55 DEBUG [example] Key: system.uname
2014/02/27 19:48:55 DEBUG [example] List: [10.8.0.1:10050 10.8.0.2:10050]
2014/02/27 19:48:55 DEBUG [example] Sending key: system.uname to 10.8.0.1:10050
2014/02/27 19:48:55 DEBUG [example] Sending key: system.uname to 10.8.0.2:10050
2014/02/27 19:48:55 DEBUG [example] Replied from 10.8.0.1:10050 in 68 msec: Linux www 3.2.0-25-generic #40-Ubuntu SMP Wed May 23 20:30:51 UTC 2012 x86_64 x86_64 x86_64 GNU/Linux
2014/02/27 19:48:55 DEBUG [example] Replied from 10.8.0.2:10050 in 86 msec: Linux raspberrypi 3.6.11+ #474 PREEMPT Thu Jun 13 17:14:42 BST 2013 armv6l GNU/Linux
2014/02/27 19:48:55 DEBUG [example] Aggregated system.uname = Linux www 3.2.0-25-generic #40-Ubuntu SMP Wed May 23 20:30:51 UTC 2012 x86_64 x86_64 x86_64 GNU/Linux
Linux raspberrypi 3.6.11+ #474 PREEMPT Thu Jun 13 17:14:42 BST 2013 armv6l GNU/Linux

2014/02/27 19:48:55 DEBUG [example] Closing connection: 127.0.0.1:65041
```

## Author

Fujiwara Shunichiro <fujiwara.shunichiro@gmail.com>
