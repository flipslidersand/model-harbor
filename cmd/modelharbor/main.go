package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/flipslidersand/model-harbor/internal/provider"
	"github.com/flipslidersand/model-harbor/internal/router"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "modelharbor",
		Short: "Multi-model AI routing gateway",
	}

	var addr string
	serve := &cobra.Command{
		Use:   "serve",
		Short: "Start the model proxy server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(addr)
		},
	}
	serve.Flags().StringVar(&addr, "addr", ":8080", "listen address")
	root.AddCommand(serve)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func runServe(addr string) error {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	harborKey := os.Getenv("MODEL_HARBOR_API_KEY")

	oai := provider.NewOpenAI(openaiKey, "")
	ant := provider.NewAnthropic(anthropicKey, "")

	rt := router.New(oai, ant)
	handler := router.BearerAuth(harborKey, rt)

	if harborKey == "" {
		fmt.Println("warning: MODEL_HARBOR_API_KEY not set — running without authentication")
	}
	fmt.Printf("modelharbor listening on %s\n", addr)
	return http.ListenAndServe(addr, handler)
}
