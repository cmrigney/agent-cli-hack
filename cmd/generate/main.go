package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cmrigney/agent-cli-hack/internal/config"
	"github.com/docker/mcp-registry/pkg/catalog"
	"gopkg.in/yaml.v3"
)

func main() {
	servers, err := readCatalog()
	if err != nil {
		log.Fatal(err)
	}

	for _, server := range servers {
		fmt.Println("Generating agent for", server.Name)
		agentConfig := convertToAgentConfig(server)

		if err := os.MkdirAll(filepath.Join("agents", server.Name), 0755); err != nil {
			log.Fatal(err)
		}

		b, err := yaml.Marshal(agentConfig)
		if err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile(filepath.Join("agents", server.Name, "agent.yaml"), b, 0644); err != nil {
			log.Fatal(err)
		}
	}
}

func readCatalog() (catalog.TileList, error) {
	resp, err := http.Get("https://desktop.docker.com/mcp/catalog/v2/catalog.yaml")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var catalog catalog.TopLevel
	if err := yaml.NewDecoder(resp.Body).Decode(&catalog); err != nil {
		return nil, err
	}

	return catalog.Registry, nil
}

func convertToAgentConfig(server catalog.TileEntry) config.AgentConfig {
	return config.AgentConfig{
		Name:        server.Name + "_agent",
		Model:       "{{.Model}}",
		Description: "An agent that wraps an mcp server. The MCP server description is: " + server.Tile.Description,
		Instruction: "You are an agent that specializes in using the " + server.Name + " mcp server. Its description is as follows:\n```" +
			server.Tile.Description +
			"```\nYour purpose is to fulfill requests using the mcp server tools provided.",
		Toolsets: []config.Toolset{
			{
				Type:    "mcp",
				Command: "docker",
				Args: []string{
					"mcp", "gateway", "run",
					"--servers={{.Root}}/" + server.Name,
					"--config={{.Root}}/gateway-config.yaml",
					"--secrets={{.Root}}/.gateway-secrets",
				},
			},
		},
	}
}
