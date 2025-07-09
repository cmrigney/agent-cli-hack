package agentrunner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cmrigney/agent-cli-hack/internal/config"
	"github.com/cmrigney/agent-cli-hack/internal/registry"
	"github.com/cmrigney/agent-cli-hack/internal/user"
	"gopkg.in/yaml.v3"
)

type Options struct {
	Model         string
	UseLocalFiles bool
	CagentPath    string
	Web           bool
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

	tempDir, err := os.MkdirTemp("", "agent-cli-hack")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	templateModel := config.TemplateModel{
		Model: r.options.Model,
		Root:  tempDir,
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

	agentFile, err := copyAgentFiles(tempDir, config)
	if err != nil {
		return err
	}

	if r.options.Web {
		return r.runCagentWeb(ctx, agentFile)
	}

	return r.runCagent(ctx, agentFile)
}

func (r *AgentRunner) runCagent(ctx context.Context, agentFile string) error {
	cmd := exec.CommandContext(ctx, r.options.CagentPath, "run", agentFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (r *AgentRunner) runCagentWeb(ctx context.Context, agentFile string) error {
	cmd := exec.CommandContext(ctx, r.options.CagentPath, "web", "-d", filepath.Dir(agentFile))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func copyAgentFiles(rootDir string, config config.Config) (string, error) {
	yaml, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(rootDir, "agents.yaml")
	if err := os.WriteFile(configPath, yaml, 0644); err != nil {
		return "", err
	}

	mcpRoot, err := mcpToolkitRoot()
	if err != nil {
		return "", err
	}

	if err := copyFile(filepath.Join(mcpRoot, "config.yaml"), filepath.Join(rootDir, "gateway-config.yaml")); err != nil {
		return "", err
	}

	return filepath.Join(rootDir, "agents.yaml"), nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func mcpToolkitRoot() (string, error) {
	homeDir, err := user.HomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".docker", "mcp"), nil
}
