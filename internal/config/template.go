package config

import (
	"bytes"
	"io"
	"text/template"

	"gopkg.in/yaml.v3"
)

type TemplateModel struct {
	Model         string
	Root          string
	SubAgentsList string
	SubAgents     string
}

func DecodeTemplate(r io.Reader, model TemplateModel) (AgentConfig, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return AgentConfig{}, err
	}

	tmpl, err := template.New("agent").Parse(string(b))
	if err != nil {
		return AgentConfig{}, err
	}

	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, model); err != nil {
		return AgentConfig{}, err
	}

	var agentConfig AgentConfig
	if err := yaml.NewDecoder(bytes.NewReader(buff.Bytes())).Decode(&agentConfig); err != nil {
		return AgentConfig{}, err
	}

	return agentConfig, nil
}
