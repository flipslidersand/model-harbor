package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/flipslidersand/model-harbor/internal/httpapi"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "modelharbor",
		Short: "Multi-model AI routing gateway",
	}

	root.AddCommand(serveCmd())
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func serveCmd() *cobra.Command {
	var addr string
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the OpenAI-compatible HTTP gateway",
		RunE: func(_ *cobra.Command, _ []string) error {
			keys := httpapi.LoadKeysFromEnv()
			if len(keys) == 0 {
				log.Printf("WARNING: no API keys configured (MODELHARBOR_API_KEYS); "+
					"/v1/chat/completions will reject all requests. Set %s to enable access.",
					"MODELHARBOR_API_KEYS")
			}
			srv := &http.Server{
				Addr:              addr,
				Handler:           httpapi.NewMux(keys),
				ReadHeaderTimeout: 10 * time.Second,
			}
			log.Printf("modelharbor listening on %s", addr)
			return srv.ListenAndServe()
		},
	}
	cmd.Flags().StringVar(&addr, "addr", ":8080", "listen address")
	return cmd
}
