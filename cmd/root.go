package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "generate",
	Short: "CLI para gerar módulos de API em Go",
	Long:  `Uma ferramenta para automatizar a criação de módulos com model, repository, service e controller.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
