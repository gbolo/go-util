# Configuration Management written in go

**This is a not a serious project, just done for educational value. DO NOT USE THIS!**

This is a simple CM solution written in go. It is `agent` based (agents expose an http REST API).
The `HQ` reads in a simple `DSL` in yaml which then gets translated to API calls to specified agents.

## Sample Configs
 - agent: `testdata/agent.yaml`
 - HQ: `testdata/dsl.yaml`

## Instructions
 // TODO

## Design

### Agent
the agent is a simple daemon binary that is running on the target host.
It exposes a simple REST API with only one endpoint: `api/v1/task`.
When a POST request is made to this endpoint, it will attempt to execute the corresponding "module"

### HQ
hq is a simple binary cli that reads in a `DSL` file on disk.
hq will attempt to make http POSTs to each target host for each defined task for "modules" which it understands


## Features (and lack of features)
 // TODO

## Q&A

**Why agent based?**
since go is a compiled language, there is no corresponding runtime available on the target host.
Without an agent, the hq would simply be executing commands over ssh (might as well write this in bash then).
Therefore, I concluded that I could only do "interesting" things if the target host runs an agent.

**Why Go?**
it's my most familiar language at the moment.
Although i realize that it probably wasn't the best choice for this :)
