package config

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func GetAgentConfig(ctx context.Context, agent string, model TemplateModel, useLocal bool) (AgentConfig, error) {
	var reader io.Reader

	if useLocal {
		file, err := os.Open(filepath.Join("agents", agent, "agent.yaml"))
		if err != nil {
			return AgentConfig{}, err
		}
		defer file.Close()
		reader = file
	} else {
		resp, err := makeRequest(ctx, "https://raw.githubusercontent.com/cmrigney/agent-cli-hack/refs/heads/main/agents/"+agent+"/agent.yaml")
		if err != nil {
			return AgentConfig{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return AgentConfig{}, fmt.Errorf("agent %s not found in the agent registry", agent)
		}
		reader = resp.Body
	}

	return DecodeTemplate(reader, model)
}

func GetCoordinatorConfig(ctx context.Context, name string, model TemplateModel, useLocal bool) (AgentConfig, error) {
	var reader io.Reader

	if useLocal {
		file, err := os.Open(filepath.Join("agent-coordinators", name+".yaml"))
		if err != nil {
			return AgentConfig{}, err
		}
		defer file.Close()
		reader = file
	} else {
		resp, err := makeRequest(ctx, "https://raw.githubusercontent.com/cmrigney/agent-cli-hack/refs/heads/main/agent-coordinators/"+name+".yaml")
		if err != nil {
			return AgentConfig{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return AgentConfig{}, fmt.Errorf("agent %s not found in the agent coordinators", name)
		}
		reader = resp.Body
	}

	return DecodeTemplate(reader, model)
}

func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: 20 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
