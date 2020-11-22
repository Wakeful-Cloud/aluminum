package steam

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/shirou/gopsutil/process"
	"github.com/wakeful-cloud/vdf"
)

//open the URI with the OS-level default URI handler
func open(uri string) error {
	//Get rundll32.exe location
	rundll32 := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")

	//Run the command
	cmd := exec.Command(rundll32, "url.dll,FileProtocolHandler", uri)
	err := cmd.Start()

	return err
}

//AddGame adds a game to Steam
func AddGame(name string, target string) error {
	return Update(func(shortcuts *vdf.Map) error {
		//Get largest index
		largestIndex := uint32(1)
		for k := range (*shortcuts)["shortcuts"].(vdf.Map) {
			//Get current index
			rawIndex, err := strconv.ParseUint(k, 10, 32)

			if err != nil {
				return err
			}

			//Convert to uint32
			index := uint32(rawIndex)

			//Update index
			if index > largestIndex {
				largestIndex = index
			}
		}

		//Generate new index
		index := fmt.Sprintf("%v", largestIndex+1)

		//Add game (https://developer.valvesoftware.com/wiki/Add_Non-Steam_Game)
		(*shortcuts)["shortcuts"].(vdf.Map)[index] = vdf.Map{
			"AllowDesktopConfig": uint32(1),
			"AllowOverlay":       uint32(1),
			"AppName":            name,
			"Devkit":             uint32(0),
			"DevkitGameID":       "",
			"Exe":                fmt.Sprintf("\"%s\"", target),
			"icon":               "",
			"IsHidden":           uint32(0),
			"LastPlayTime":       uint32(0),
			"LaunchOptions":      "",
			"OpenVR":             uint32(0),
			"ShortcutPath":       "",
			"StartDir":           fmt.Sprintf("\"%s\"", filepath.Dir(target)),
			"tags":               vdf.Map{},
		}

		return nil
	})
}

//CalculateID calculates the Steam ID for something
func CalculateID(name string, target string) string {
	//Calculate checksum
	checksum := crc32.ChecksumIEEE([]byte(fmt.Sprintf("\"%s\"", target) + name))

	//Math
	a := checksum | 0x80000000
	b := uint64(a) << 32
	c := b | 0x02000000

	return fmt.Sprintf("%v", c)
}

//CheckGame checks wether or not a game is registered with Steam
func CheckGame(name string, target string) (bool, error) {
	//Get userdata directory
	userdata := filepath.Join(os.Getenv("ProgramFiles(x86)"), "Steam", "userdata")

	//Get all userdata subdirectories
	subdirectories, err := ioutil.ReadDir(userdata)

	if err != nil {
		return false, err
	}

	//Get shortcut path
	shortcutPath := filepath.Join(userdata, subdirectories[0].Name(), "config", "shortcuts.vdf")

	//Read the shortcuts
	bytes, err := ioutil.ReadFile(shortcutPath)

	if err != nil {
		return false, err
	}

	//Parse steam shortcuts
	shortcuts, err := vdf.ReadVdf(bytes)

	if err != nil {
		return false, err
	}

	for _, v := range shortcuts["shortcuts"].(vdf.Map) {
		if v.(vdf.Map)["Exe"].(string) == fmt.Sprintf("\"%s\"", target) &&
			v.(vdf.Map)["AppName"].(string) == name {
			return true, nil
		}
	}

	return false, nil
}

//Open a game with Steam
func Open(id string) error {
	//Generate the URI
	uri := fmt.Sprintf("steam://rungameid/%s", id)

	return open(uri)
}

//Stop steam
func Stop() error {
	//Get all processes
	processes, err := process.Processes()

	if err != nil {
		panic(err)
	}

	//Stop Steam
	for _, v := range processes {
		name, err := v.Name()

		if err == nil && name == "steam.exe" {
			err = v.Kill()

			if err != nil {
				return err
			}
		}
	}

	return nil
}

//RemoveGame removes a game from Steam
func RemoveGame(name string, target string) error {
	return Update(func(shortcuts *vdf.Map) error {
		for k, v := range (*shortcuts)["shortcuts"].(vdf.Map) {
			if v.(vdf.Map)["Exe"].(string) == fmt.Sprintf("\"%s\"", target) &&
				v.(vdf.Map)["AppName"].(string) == name {
				delete((*shortcuts)["shortcuts"].(vdf.Map), k)
			}
		}

		return nil
	})
}

//UpdateArguments updates a game's Steam launch arguments
func UpdateArguments(name string, target string, args string) error {
	return Update(func(shortcuts *vdf.Map) error {
		//Find appropriate shortcut and modify launch arguments
		for _, v := range (*shortcuts)["shortcuts"].(vdf.Map) {
			if v.(vdf.Map)["Exe"].(string) == fmt.Sprintf("\"%s\"", target) &&
				v.(vdf.Map)["AppName"].(string) == name {
				//Modify launcher options
				v.(vdf.Map)["LaunchOptions"] = args
			}
		}

		return nil
	})
}

//Update Steam shortcuts
func Update(updater func(*vdf.Map) error) error {
	//Get userdata directory
	userdata := filepath.Join(os.Getenv("ProgramFiles(x86)"), "Steam", "userdata")

	//Get all userdata subdirectories
	subdirectories, err := ioutil.ReadDir(userdata)

	if err != nil {
		return err
	}

	//Get shortcut path
	shortcutPath := filepath.Join(userdata, subdirectories[0].Name(), "config", "shortcuts.vdf")

	//Read the shortcuts
	bytes, err := ioutil.ReadFile(shortcutPath)

	if err != nil {
		return err
	}

	//Parse steam shortcuts
	shortcuts, err := vdf.ReadVdf(bytes)

	if err != nil {
		return err
	}

	//Invoke updater
	err = updater(&shortcuts)

	if err != nil {
		return err
	}

	//Convert back to VDF
	bytes, err = vdf.WriteVdf(shortcuts)

	if err != nil {
		return err
	}

	//Save new VDF
	err = ioutil.WriteFile(shortcutPath, bytes, 0644)

	if err != nil {
		return err
	}

	return nil
}
