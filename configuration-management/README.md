# Configuration Management written in go

**This is a not a serious project, just done for educational value. DO NOT USE THIS!**

This is a simple CM solution written in go. It is `agent` based (agents expose an http REST API).
The `HQ` reads in a simple `DSL` in yaml which then gets translated to API calls to specified agents.

## Sample Configs
 - agent: `testdata/agent.yaml`
 - HQ: `testdata/dsl.yaml`

## Test Instructions
**Docker is required to follow these test instructions.**
 1. build the docker image: `./build-docker.sh`
 2. run some agents: `./run-agent.sh 18001 && ./run-agent.sh 18002`
 3. (optional) modify the hq dsl file: `vi testdata/dsl.yaml`
 4. run the hq: `./run-hq.sh`
 5. (optional) view agent logs `docker logs agent18001`

## Design

### Agent
the agent is a simple daemon binary that is running on the target host.
It exposes a simple REST API with only one endpoint: `api/v1/task`.
When a POST request is made to this endpoint, it will attempt to execute the corresponding "module"

### HQ
hq is a simple binary cli that reads in a `DSL` file on disk.
hq will attempt to make http POSTs to each target host for each defined task for "modules" which it understands


## Modules
the following "modules" are implemented.
**NOTE: these modules are basic, they will ALWAYS execute the task regardless of the current state of the resource**

### directory
```
# create/delete a directory
- module: directory
  # path to the directory
  name: /tmp/test
  # valid states: present, absent
  state: absent
```

### file
```
# create/delete a file
- module: file
  # path to the file
  name: /tmp/test/f1.txt
  # valid states: present, absent
  state: absent
  # the content of the file (if state is present)
  content: |
      some test content
      for file f1.txt
```

### apt (package management)
**This module requires that `apt-get` be installed on target host**
```
# install/remove a package via apt-get
- module: apt
  # the name of the package
  name: htop
  # valid states: present, absent
  state: absent

# setting state to update will update the apt-get cache
- module: apt
  state: update
```

### service
**This module requires that `systemd` is the init system for the target host**
```
# start/stop/restart a systemd service
- module: service
  # the name of the service
  name: ntp
  # valid states: start, stop, restart
  state: start
```

### shellcmd (run a shell command)
**This module will execute a command on a new `/bin/sh` shell.**
The command output is ignored. Success depends on exit code.
It is expected that the command NOT be interactive and terminates prior to the http server timeout settings :)
```
# run a shell command
- module: shellcmd
  # the shell command to run
  name: ls -alt > ls.output
```

## Q&A

**Why agent based?**
since go is a compiled language, there is no corresponding runtime available on the target host.
Without an agent, the hq would simply be executing commands over ssh (might as well write this in bash then).
Therefore, I concluded that I could only do "interesting" things if the target host runs an agent.

**Why Go?**
it's my most familiar language at the moment.
Although i realize that it probably wasn't the best choice for this :)

## TODO
 - add authentication to agent http server (via auth http header and/or TLS mutual auth)
 - make modules understand current state of resource and improve result returned
 - add support for other OS (for package module) like alpine and/or centOS
 - submit each task to every target host in parallel
