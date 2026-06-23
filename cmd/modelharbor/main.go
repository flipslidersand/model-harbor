package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "modelharbor",
		Short: "Multi-model AI routing gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("modelharbor — not yet implemented")
			return nil
		},
	}
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
