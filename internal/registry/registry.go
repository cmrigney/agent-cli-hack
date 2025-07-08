package registry

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cmrigney/agent-cli-hack/internal/config"
	"gopkg.in/yaml.v3"
)

// https://raw.githubusercontent.com/cmrigney/agent-cli-hack/refs/heads/main/agents/astra-db/agent.yaml

func EnableServer(ctx context.Context, name string, useLocal bool) error {
	if err := ensureRegistryFile(); err != nil {
		return err
	}

	// Verify that the agent exists
	if _, err := config.GetAgentConfig(ctx, name, config.TemplateModel{
		// Dummy data
		Model: "test",
		Root:  "testpath",
	}, useLocal); err != nil {
		return err
	}

	registry, err := read()
	if err != nil {
		return err
	}

	registry.Agents[name] = true

	if err := write(registry); err != nil {
		return err
	}

	return nil
}

func DisableServer(name string) error {
	if err := ensureRegistryFile(); err != nil {
		return err
	}

	registry, err := read()
	if err != nil {
		return err
	}

	delete(registry.Agents, name)

	if err := write(registry); err != nil {
		return err
	}

	return nil
}

func GetRegisteredAgents(ctx context.Context, model config.TemplateModel, useLocal bool) ([]config.AgentConfig, error) {
	if err := ensureRegistryFile(); err != nil {
		return nil, err
	}

	registry, err := read()
	if err != nil {
		return nil, err
	}

	agents := make([]config.AgentConfig, 0, len(registry.Agents))
	for name := range registry.Agents {
		agent, err := config.GetAgentConfig(ctx, name, model, useLocal)
		if err != nil {
			return nil, err
		}

		agents = append(agents, agent)
	}

	return agents, nil
}

func ensureRegistryFile() error {
	err := os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".docker", "agents"), 0755)
	if err != nil {
		return err
	}

	_, err = os.Stat(filepath.Join(os.Getenv("HOME"), ".docker", "agents", "registry.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			registry := Registry{
				Agents: make(map[string]bool),
			}
			b, err := yaml.Marshal(registry)
			if err != nil {
				return err
			}
			return os.WriteFile(filepath.Join(os.Getenv("HOME"), ".docker", "agents", "registry.yaml"), b, 0644)
		}
	}

	return nil
}

func read() (Registry, error) {
	// TODO: only works on Mac and with DD
	b, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".docker", "agents", "registry.yaml"))
	if err != nil {
		return Registry{}, err
	}

	var registry Registry
	if err := yaml.Unmarshal(b, &registry); err != nil {
		return Registry{}, err
	}

	return registry, nil
}

func write(registry Registry) error {
	b, err := yaml.Marshal(registry)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(os.Getenv("HOME"), ".docker", "agents", "registry.yaml"), b, 0644)
}
