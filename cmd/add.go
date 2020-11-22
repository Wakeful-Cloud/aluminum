package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wakeful-cloud/aluminum/steam"
)

//addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a game",
	Long: `Aluminum will rename your game and generate a mock game with the same name.
When you normally start your game (Such as through another game launcher),
the other launcher will run the mock game which will start Aluminum which will handle
everything else.

If you accidentally run this command or want to remove Aluminum, you can run the "remove"
command to undo anything done by running this command in the first place.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("This only takes the location of the game file")
		}

		//Check file statistics
		_, err := os.Stat(args[0])

		if err == nil {
			return nil
		} else if os.IsNotExist(err) {
			return fmt.Errorf("Game file: %s doesn't exist", args[0])
		} else {
			return fmt.Errorf("Unknown game file %s state %s", args[0], err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		//Alias target
		target := args[0]

		//Get path to mock game
		_, filename, _, _ := runtime.Caller(0)

		mock := path.Join(filepath.Dir(filepath.Dir(filename)), "mock/mock.exe")

		//Make sure it exists
		_, err := os.Stat(mock)

		if err == nil {
			//Parse game target
			dir := filepath.Dir(target)
			extension := filepath.Ext(target)
			name := strings.TrimSuffix(filepath.Base(target), extension)

			//Generate new target
			newTarget := filepath.Join(dir, fmt.Sprintf("%sAluminum%s", name, extension))

			//Add to database
			db.Put([]byte(name), []byte(target))

			//Rename
			err := os.Rename(target, newTarget)

			if err != nil {
				panic(err)
			}
			//Link to mock game
			err = os.Link(mock, target)

			if err != nil {
				panic(err)
			}

			//Generate game config
			game := game{
				Name:   fmt.Sprintf("%s (Aluminum)", name),
				Target: newTarget,
			}

			//Marshal
			bytes, err := json.Marshal(game)

			if err != nil {
				panic(err)
			}

			//Generate config path
			configPath := filepath.Join(dir, "aluminum-config.json")

			//Write
			err = ioutil.WriteFile(configPath, bytes, 0644)

			if err != nil {
				panic(err)
			}

			//Stop Steam
			err = steam.Stop()

			if err != nil {
				panic(err)
			}

			//Add the game to Steam
			err = steam.AddGame(fmt.Sprintf("%s (Aluminum)", name), newTarget)

			if err != nil {
				panic(err)
			}

			fmt.Printf("Added game %s (At %s)\n", name, target)
		} else if os.IsNotExist(err) {
			fmt.Printf("You need to build the mock game first! (See the installations instructions!)")
			os.Exit(1)
		} else {
			fmt.Printf("Unknown mock game file %s state %s", target, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
