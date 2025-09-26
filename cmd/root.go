package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command {
	Use: "pass_manager",
	Short: "password manager",
}

func Execute() error {
	return rootCmd.Execute() // avvia e controlla se è stato inserito un comando valido
}