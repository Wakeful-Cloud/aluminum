package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "aluminum",
	Short: "Use non-Steam games with Steam without compromise",
	Long: `Aluminum is a tool that generates drop-in replacements
for your game that launch the game through Steam instead of
directly launching it.`,
}

//Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
