package cmd

import (
	"encoding/json"
	"fmt"
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
				//Wether or not the game was OK
				ok := true

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
				_, err = os.Stat(newTarget)

				if err != nil && os.IsNotExist(err) {
					ok = false

					fmt.Printf("Game at %s is missing\n", target)

					if !dry {
						//Rename
						err = os.Rename(target, newTarget)

						//Link to mock game
						err = os.Link(mock, target)

						if err != nil {
							panic(err)
						}

						fmt.Printf("Fixed game at %s\n", target)
					}
				} else if err != nil {
					return fmt.Errorf("Unknown game %s state %s", target, err)
				}

				//Verify config file exists
				configTarget := filepath.Join(dir, "aluminum-config.json")
				_, err = os.Stat(configTarget)

				if err != nil && os.IsNotExist(err) {
					ok = false

					fmt.Printf("Config file at %s is missing\n", configTarget)

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
						err = os.WriteFile(configTarget, bytes, 0644)

						if err != nil {
							panic(err)
						}

						fmt.Printf("Fixed config file at %s\n", configTarget)
					}
				} else if err != nil {
					return fmt.Errorf("Unknown config file %s state %s", configTarget, err)
				}

				//Check if game is in Steam
				hasGame, err := steam.CheckGame(fmt.Sprintf("%s (Aluminum)", name), newTarget)

				if err != nil {
					panic(err)
				}

				//If the game isn't registered with Steam, add it
				if !hasGame {
					ok = false

					fmt.Printf("Steam entry for %s is missing\n", newTarget)

					if !dry {
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

						fmt.Printf("Added Steam entry for %s\n", newTarget)
					}
				}

				if ok {
					fmt.Printf("Verified %s (At %s) - found no problems\n", name, target)
				} else {
					fmt.Printf("Verified %s (At %s) - found problems \n", name, target)
				}

				return nil
			})

			if err != nil {
				panic(err)
			}
		} else if os.IsNotExist(err) {
			fmt.Println("You need to build the mock game first! (See the installations instructions!)")
			os.Exit(1)
		} else {
			fmt.Printf("Unknown mock game file %s state %s\n", mock, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().BoolVar(&dry, "dry", false, "Dry run (Only lists missing mock games but doesn't fix them)")
}
