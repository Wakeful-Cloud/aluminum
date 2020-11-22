package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wakeful-cloud/aluminum/steam"
)

var all bool

//removeGame removes a game
func removeGame(name string, target string) error {
	//Parse game target
	dir := filepath.Dir(target)
	extension := filepath.Ext(target)

	//Remove mock game
	err := os.Remove(target)

	if err != nil {
		return err
	}

	//Generate config path
	configPath := filepath.Join(dir, "aluminum-config.json")

	//Remove config
	err = os.Remove(configPath)

	if err != nil {
		return err
	}

	//Generate log path
	logPath := filepath.Join(dir, "aluminum-log.txt")

	//If the log file exists, remove it
	_, err = os.Stat(logPath)

	if err == nil {
		//Remove config
		err = os.Remove(logPath)

		if err != nil {
			return err
		}
	}

	//Generate old target
	oldTarget := filepath.Join(dir, fmt.Sprintf("%sAluminum%s", name, extension))

	//Rename
	err = os.Rename(oldTarget, target)

	//Stop Steam
	err = steam.Stop()

	if err != nil {
		panic(err)
	}

	//Remove the game from Steam
	err = steam.RemoveGame(fmt.Sprintf("%s (Aluminum)", name), oldTarget)

	if err != nil {
		panic(err)
	}

	return nil
}

//removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a game",
	Long: `Aluminum will delete the mock game and rename your game back to its original name.

This command effectively undoes what the "add" command does.`,
	Args: func(cmd *cobra.Command, args []string) error {
		//All switch ignores all arguments
		if all {
			return nil
		}

		if len(args) != 1 {
			return errors.New("This only takes the game name")
		}

		//Ensure game exists in the database
		if db.Has([]byte(args[0])) {
			return nil
		}

		return fmt.Errorf("Game: %s not registered with Aluminum", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		if all {
			//Iterate over games
			err := db.Fold(func(rawName []byte) error {
				//Get raw target
				rawTarget, err := db.Get(rawName)

				if err != nil {
					return err
				}

				//Convert to string
				name := string(rawName)
				target := string(rawTarget)

				//Remove the game
				err = removeGame(name, target)

				if err != nil {
					return err
				}

				fmt.Printf("Removed game %s (At %s)\n", name, target)

				return nil
			})

			if err != nil {
				panic(err)
			}

			//Remove from database
			err = db.DeleteAll()

			if err != nil {
				panic(err)
			}
		} else {
			//Alias name
			name := args[0]

			//Get game target
			rawTarget, err := db.Get([]byte(name))

			if err != nil {
				panic(err)
			}

			//Convert to string
			target := string(rawTarget)

			//Remove the game
			removeGame(name, target)

			//Remove from database
			err = db.Delete([]byte(name))

			if err != nil {
				panic(err)
			}

			fmt.Printf("Removed game %s (At %s)\n", name, target)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolVar(&all, "all", false, "Remove all mock games configured with Aluminum")
}
