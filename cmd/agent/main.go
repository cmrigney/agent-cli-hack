package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		return rootCommand(ctx)
	}, manager.Metadata{
		SchemaVersion:    "0.1.0",
		Vendor:           "Cody Rigney",
		Version:          Version,
		ShortDescription: "Docker Agent Plugin",
	})
}
