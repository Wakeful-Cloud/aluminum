package cmd

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/prologic/bitcask"
)

var db *bitcask.Bitcask

func init() {
	//Get home directory
	home, err := homedir.Dir()

	if err != nil {
		panic(err)
	}

	//Generate database path
	dbPath := filepath.Join(home, "aluminum")

	//Open database
	db, err = bitcask.Open(dbPath)

	if err != nil {
		panic(err)
	}
}
