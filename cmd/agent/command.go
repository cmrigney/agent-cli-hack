package main

import (
	"context"
	"fmt"

	"github.com/cmrigney/agent-cli-hack/internal/agentrunner"
	"github.com/cmrigney/agent-cli-hack/internal/registry"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/spf13/cobra"
)

const Version = "HEAD"

func rootCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "agent [OPTIONS]",
		TraverseChildren: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: false,
			HiddenDefaultCmd:  true,
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetContext(ctx)
			return plugin.PersistentPreRunE(cmd, args)
		},
		Version: Version,
	}

	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.Flags().BoolP("version", "v", false, "Print version information and quit")

	_ = cmd.RegisterFlagCompletionFunc("agent", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"--help"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(versionCommand())
	cmd.AddCommand(runCommand())
	cmd.AddCommand(enableCommand())
	cmd.AddCommand(disableCommand())
	cmd.AddCommand(listCommand())

	return cmd
}

func enableCommand() *cobra.Command {
	var useLocal bool

	cmd := &cobra.Command{
		Use:   "enable",
		Args:  cobra.ExactArgs(1),
		Short: "Enable an agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := registry.EnableAgent(cmd.Context(), args[0], useLocal); err != nil {
				return err
			}

			fmt.Println("Enabled agent", args[0])
			return nil
		},
	}

	cmd.Flags().BoolVar(&useLocal, "use-local", false, "Use local files instead of fetching from the registry")

	return cmd
}

func disableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Args:  cobra.ExactArgs(1),
		Short: "Disable an agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := registry.DisableAgent(args[0]); err != nil {
				return err
			}

			fmt.Println("Disabled agent", args[0])
			return nil
		},
	}

	return cmd
}

func listCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Args:  cobra.NoArgs,
		Short: "List all enabled agents",
		RunE: func(cmd *cobra.Command, args []string) error {
			agents, err := registry.ListAgents()
			if err != nil {
				return err
			}

			for _, agent := range agents {
				fmt.Println(agent)
			}

			return nil
		},
	}

	return cmd
}

func runCommand() *cobra.Command {
	options := agentrunner.Options{
		Model:          "gpt-4o",
		UseLocalFiles:  false,
		CagentPath:     "../cagent/bin/cagent", // assume it's one level up
		Web:            false,
		Debug:          false,
		Think:          true,
		ThinkSubAgents: true,
		Todo:           false,
	}

	cmd := &cobra.Command{
		Use:   "run",
		Args:  cobra.NoArgs,
		Short: "Interact with the agent team",
		RunE: func(cmd *cobra.Command, args []string) error {
			return agentrunner.NewAgentRunner(options).Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&options.Model, "model", "gpt-4o", "Model to use for the agents")
	cmd.Flags().BoolVar(&options.UseLocalFiles, "use-local", false, "Use local files instead of fetching from the registry")
	cmd.Flags().StringVar(&options.CagentPath, "cagent", "../cagent/bin/cagent", "Path to the cagent binary")
	cmd.Flags().BoolVar(&options.Web, "web", false, "Run the cagent web interface instead of the CLI")
	cmd.Flags().BoolVar(&options.Debug, "debug", false, "Debug mode provides verbose output")
	cmd.Flags().BoolVar(&options.Think, "think", true, "Enable thinking for the coordinator")
	cmd.Flags().BoolVar(&options.ThinkSubAgents, "think-subagents", true, "Enable thinking for subagents")
	cmd.Flags().BoolVar(&options.Todo, "todo", false, "Enable todo list for the coordinator")

	return cmd
}

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Short: "Show the version information",
		Use:   "version",
		Args:  cobra.ExactArgs(0),
		// Deactivate PersistentPreRun for this command only
		// We don't want to check if Docker Desktop is running.
		PersistentPreRun: func(*cobra.Command, []string) {},
		Run: func(*cobra.Command, []string) {
			fmt.Println(Version)
		},
	}
}
