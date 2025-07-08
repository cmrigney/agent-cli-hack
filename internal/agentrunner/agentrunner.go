package agentrunner

import (
	"context"
	"fmt"

	"github.com/cmrigney/agent-cli-hack/internal/config"
	"github.com/cmrigney/agent-cli-hack/internal/registry"
	"gopkg.in/yaml.v3"
)

type Options struct {
	Model         string
	UseLocalFiles bool
}

func (o *Options) Validate() error {
	switch o.Model {
	case "gpt-4o":
		break
	default:
		return fmt.Errorf("model %s is not supported", o.Model)
	}

	return nil
}

type AgentRunner struct {
	options Options
}

func NewAgentRunner(options Options) *AgentRunner {
	return &AgentRunner{
		options: options,
	}
}

func (r *AgentRunner) Run(ctx context.Context) error {
	if err := r.options.Validate(); err != nil {
		return err
	}

	templateModel := config.TemplateModel{
		Model: r.options.Model,
		Root:  "testpath",
	}

	agents, err := registry.GetRegisteredAgents(ctx, templateModel, r.options.UseLocalFiles)
	if err != nil {
		return err
	}

	subAgentsListString := ""
	subAgentsString := "["
	first := true
	for _, agent := range agents {
		if !first {
			subAgentsString += ", "
			subAgentsListString += "\n  "
		}
		subAgentsString += fmt.Sprintf(`"%s"`, agent.Name)
		subAgentsListString += fmt.Sprintf(`- %s: %s`, agent.Name, agent.Description)
		first = false
	}
	subAgentsString += "]"

	templateModel.SubAgentsList = subAgentsListString
	templateModel.SubAgents = subAgentsString

	coordinator, err := config.GetCoordinatorConfig(ctx, "devex", templateModel, r.options.UseLocalFiles)
	if err != nil {
		return err
	}

	allAgents := make(map[string]config.AgentConfig)
	allAgents["root"] = coordinator
	for _, agent := range agents {
		allAgents[agent.Name] = agent
	}

	config := config.Config{
		Agents: allAgents,
		Models: map[string]config.ModelConfig{
			// TODO add more models
			r.options.Model: {
				Type:  "openai",
				Model: r.options.Model,
			},
		},
	}

	yaml, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	fmt.Println(string(yaml))
	return nil
}
