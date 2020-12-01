# Aluminum
![Status Shield](https://img.shields.io/badge/status-beta-yellow?style=for-the-badge)

Use non-Steam games with Steam without compromise.

## Install
1. Install [Go](https://golang.org/doc/install)
2. Download and build Aluminum: `go get github.com/wakeful-cloud/aluminum`
3. **Build mock game**:
```powershell
cd $env:GOPATH/src/github.com/wakeful-cloud/aluminum/mock; go build
```

## Uninstall
1. Remove all games: `aluminum remove --all`
2. Nuke the database: `aluminum nuke`
3. Remove the Go package:
```powershell
rm $env:GOPATH/src/github.com/wakeful-cloud/aluminum -r
rm $env:GOPATH/src/github.com/wakeful-cloud/aluminum -r
```
4. No restart is necessary 😎

## Usage
* Add games: `aluminum add "path/to/game/exe"`
* Verify games (After updating): `aluminum verify`
* List games registered with Aluminum: `aluminum list`
* Remove games: `aluminum remove [name]`

## Limitations
* Windows only

## Troubleshooting
1. Verify all Aluminum game integrations: `aluminum verify`
2. Manually inspect the game install location
    1. There should be the original game called `[name]Aluminum.[extension]`
    2. There should be the mock game `[name].[extension]`
    3. There should be a mock game config file `aluminum-config.json`
    4. There should be a log file `aluminum-log.txt`
3. Nuke the database: `aluminum nuke`

## How it works
When you add a game to Aluminum, it renames your game to `[name]Aluminum.exe`. It then symbolically
links the mock game with the same name as your original game (`[name].exe`). When you start your game
via a launcher (Other than Steam), it will launch the mock game which intercepts the CLI arguments,
updates the Steam game's launch arguments, then starts the game via Steam (You can see all of this
happening in the `aluminum-log.txt`).
