package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all games configured with Aluminum",
	Long:  `This command simply lists all games configured with Aluminum.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Iterate over games
		err := db.Fold(func(name []byte) error {
			//Get raw target
			rawTarget, err := db.Get(name)

			if err != nil {
				return err
			}

			//Convert to string
			target := string(rawTarget)

			//Print
			fmt.Printf("Name: %s Target: %s\n", name, target)

			return nil
		})

		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
