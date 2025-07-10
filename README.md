# Docker Agent CLI Hack

A hacky Docker CLI plugin for managing and running AI agent teams. Built for an internal Docker hackathon as a PoC.

## Prerequisites

- Go 1.24+ installed
- [Task](https://taskfile.dev)
- Docker Desktop (needed for plugin integration and MCP Toolkit configuration)
- cagent - ask @rumpl for access.
- Only works on MacOS right now

## Build and Install

**Prepare cagent:**

Clone cagent from the same directory that you cloned this repo and build it. The agent cli assumes cagent's binary exists at `../cagent/bin/cagent` (can be overriden with the `--cagent` flag).

**Build and install the CLI:**

```
task install
```

Now verify it's installed with:

```
docker agent version
```

## Main Features

**Manage Subagents**:

You can enable or disable subagents that can be used by the coordinator.

```
docker agent enable <some-agent>
docker agent disable <some-agent>
docker agent ls

# Or to use local agent files instead of getting them from the repo
docker agent enable --use-local <some-agent>
```

Note, all the agents thus far are auto-generated from the [Docker MCP registry](https://github.com/docker/mcp-registry). The MCP servers have been automatically converted to agent yaml files that cagent understands (with some additional templating to handle some dynamic fields). Those files live under `agents`. These are the set of agents you can enable.

**Run the Agent Team**:

Prereqs:
- If using the `gpt-4o` model (the default), you must have `OPENAI_API_KEY` set as an environment variable.
- If you are using an Anthropic model like `claude-4-sonnet-latest` or `claude-3-sonnet-latest`, you must have `ANTHROPIC_API_KEY` set as an environment variable.
- You must have configured all the tools used by the agents in the MCP Toolkit of Docker Desktop (agents correspond directly to their MCP server counterpart).

You can run the coordinator agent along with the subagents you've enabled. This command basically prepares all the configuration needed and runs cagent.

```
docker agent run

# Or to use local agent files intead of getting them from the repo
docker agent run --use-local
```

Or to use cagent's web interface on http://localhost:8080.

```
docker agent run --web
```

## Example of it in action

```
$ docker agent enable obsidian
Enabled agent obsidian

$ docker agent enable youtube_transcript
Enabled agent youtube_transcript

$ docker agent enable duckduckgo
Enabled agent duckduckgo

$ docker agent run

Enter your messages (Ctrl+C to exit):
> What can you do?
[root]: I can coordinate with my team of specialized agents to complete various tasks:

1. **Web Searches**: I can delegate tasks to the duckduckgo_agent to perform web searches and gather information.
   
2. **Note Management**: The obsidian_agent can interact with Obsidian to manage notes, by creating, reading, and updating information.

3. **Video Transcripts**: The youtube_transcript_agent retrieves transcripts from YouTube videos to provide their textual content.

If you have a specific task in mind that requires one of these services, feel free to let me know, and I’ll get started.
> 
```