package main

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/wakeful-cloud/aluminum/steam"
)

func main() {
	//Open the log
	logger, err := os.OpenFile("aluminum-log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		panic(err)
	}

	//Set log output
	log.SetOutput(logger)

	//Get config
	viper.SetConfigName("aluminum-config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	//Get game target and name
	name := viper.GetString("Name")
	target := viper.GetString("Target")

	log.Printf("Intercepting launch for name: %s, target: %s\n", name, target)

	//Stop Steam
	err = steam.Stop()

	if err != nil {
		panic(err)
	}

	log.Printf("Stopped Steam\n")

	//Get launch arguments
	args := strings.Join(os.Args, " ")

	//Update arguments
	err = steam.UpdateArguments(name, target, args)

	if err != nil {
		panic(err)
	}

	log.Print("Updated launch arguments\n")

	//Calculate ID
	id := steam.CalculateID(name, target)

	//Start Steam with game
	err = steam.Open(id)

	if err != nil {
		panic(err)
	}

	log.Printf("Started steam for game ID %s\n", id)
}
