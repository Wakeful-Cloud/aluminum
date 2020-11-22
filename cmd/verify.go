package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/wakeful-cloud/aluminum/steam"

	"github.com/spf13/cobra"
)

var dry bool

//verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify the integrity of all mock games created by Aluminum",
	Long: `If you updated a game, it likely destroyed the mock game created by Aluminum.
This command ensures that all games originally set up by Aluminum have their corrresponding
mock games correctly installed. If a mock game is missing, Aluminum will recreate it.`,
	Run: func(cmd *cobra.Command, args []string) {
		//Get path to mock game
		_, filename, _, _ := runtime.Caller(0)
		mock := path.Join(filepath.Dir(filepath.Dir(filename)), "mock/mock.exe")

		//Make sure it exists
		_, err := os.Stat(mock)

		if err == nil {
			//Iterate over games
			err := db.Fold(func(rawName []byte) error {
				//Get raw target
				rawTarget, err := db.Get(rawName)

				if err != nil {
					return err
				}

				//Convert to strings
				name := string(rawName)
				target := string(rawTarget)

				//Parse game target
				dir := filepath.Dir(target)
				extension := filepath.Ext(target)

				//Generate new target
				newTarget := filepath.Join(dir, fmt.Sprintf("%sAluminum%s", name, extension))

				//Verify mock game exists
				_, err = os.Stat(target)

				if err != nil && os.IsNotExist(err) {
					if !dry {
						//Rename
						err = os.Rename(target, newTarget)
						//Link to mock game
						err = os.Link(mock, target)

						if err != nil {
							panic(err)
						}
					}
				} else if err != nil {
					return fmt.Errorf("unknown mock game file %s state %s", target, err)
				}

				//Verify config file exists
				configTarget := filepath.Join(dir, "aluminum-config.json")
				_, err = os.Stat(configTarget)

				if err != nil && os.IsNotExist(err) {
					if !dry {
						//Marshal
						bytes, err := json.Marshal(game{
							Name:   fmt.Sprintf("%s (Aluminum)", name),
							Target: newTarget,
						})

						if err != nil {
							panic(err)
						}

						//Generate config target
						configTarget := filepath.Join(dir, "aluminum-config.json")

						//Write
						err = ioutil.WriteFile(configTarget, bytes, 0644)

						if err != nil {
							panic(err)
						}
					}
				} else if err != nil {
					return fmt.Errorf("unknown config file %s state %s", configTarget, err)
				}

				//Check if game is in Steam
				hasGame, err := steam.CheckGame(fmt.Sprintf("%s (Aluminum)", name), newTarget)

				if err != nil {
					panic(err)
				}

				//If the game isn't registered with Steam, add it
				if !hasGame {
					//Stop steam
					err = steam.Stop()

					if err != nil {
						panic(err)
					}

					//Add the game to Steam
					err = steam.AddGame(fmt.Sprintf("%s (Aluminum)", name), newTarget)

					if err != nil {
						panic(err)
					}
				}

				fmt.Printf("Verified %s (At %s)\n", name, target)

				return nil
			})

			if err != nil {
				panic(err)
			}
		} else if os.IsNotExist(err) {
			fmt.Printf("You need to build the mock game first! (See the installations instructions!)")
			os.Exit(1)
		} else {
			fmt.Printf("Unknown mock game file %s state %s", mock, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().BoolVar(&dry, "dry", false, "Dry run (Only lists corrupted mock games but doesn't fix them)")
}
