package commands

import (
	svrcmd "github.com/numaproj/numaflow/server/cmd"
	"github.com/spf13/cobra"
)

func NewServerCommand() *cobra.Command {
	var (
		insecure bool
		port     int
		baseHRef string
	)

	command := &cobra.Command{
		Use:   "server",
		Short: "Start a Numaflow server",
		Run: func(cmd *cobra.Command, args []string) {
			if !cmd.Flags().Changed("port") && insecure {
				port = 8080
			}
			svrcmd.Start(insecure, port, baseHRef)
		},
	}
	command.Flags().BoolVar(&insecure, "insecure", false, "Whether to disable TLS, defaults to false.")
	command.Flags().IntVarP(&port, "port", "p", 8443, "Port to listen on, defaults to 8443 or 8080 if insecure is set")
	command.Flags().StringVar(&baseHRef, "base-href", "/", "Base href in index.html.  Useful for when the server is running behind a reverse proxy under a path other than /")
	return command
}
