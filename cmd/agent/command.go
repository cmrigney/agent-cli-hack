package main

import (
	"context"
	"fmt"

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

	return cmd
}

func enableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Args:  cobra.ExactArgs(1),
		Short: "Enable an agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Enabling agent", args[0])
			return nil
		},
	}

	return cmd
}

func disableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Args:  cobra.ExactArgs(1),
		Short: "Disable an agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Disabling agent", args[0])
			return nil
		},
	}

	return cmd
}

func runCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Args:  cobra.NoArgs,
		Short: "Interact with the agent team",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Running agent team")
			return nil
		},
	}

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
