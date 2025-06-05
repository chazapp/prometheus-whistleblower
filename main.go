package main

import (
	"context"
	"log"
	"os"

	"github.com/chazapp/prometheus-whistleblower/server"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "run",
		Usage: "Run the Prometheus-whistleblower server",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return server.Run(cmd.Int("port"))
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Value:   8080,
				Aliases: []string{"p"},
				Usage:   "Port to serve the application to",
				Sources: cli.EnvVars("PORT"),
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
