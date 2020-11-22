package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var confirm bool

// nukeCmd represents the nuke command
var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Nuke the database",
	Long:  `Deletes the database directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		if confirm {
			//Disconnect the database
			err := db.Close()

			if err != nil {
				panic(err)
			}

			//Get home directory
			home, err := homedir.Dir()

			if err != nil {
				panic(err)
			}

			//Generate database path
			dbPath := filepath.Join(home, "aluminum")

			//Delete the database directory
			err = os.RemoveAll(dbPath)

			if err != nil {
				panic(err)
			}

			fmt.Println("Database nuked!")
		} else {
			fmt.Print(`You're attempting to nuke the database!

Nuking the database is irreversible! Aluminum will loose track of
what games you have installed and will be unable to remove game
integrations automatically (You'll have to do this manually!). If
you're absolutely sure you want to do this, please re-run this
command with the "--confirm" flag. You've been warned!
`)
		}
	},
}

func init() {
	rootCmd.AddCommand(nukeCmd)

	nukeCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm nuking the database")
}
