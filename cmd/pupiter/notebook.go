package main

import (
	"fmt"

	"github.com/evanboyle/pupiter/server"
	"github.com/spf13/cobra"
)

func NewNotebookCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "notebook",
		Short: "Launch a local notebook interface.",
		Long:  `Launch a local notebook interface.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting up notebook on port 1337...")
			fmt.Println("ctrl+c to quit")
			server.Serve()
		},
	}
}
