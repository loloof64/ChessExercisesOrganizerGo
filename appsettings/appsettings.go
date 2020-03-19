package appsettings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/wailsapp/wails"
)

// SettingsManager handles application settings
type SettingsManager struct {
	Runtime *wails.Runtime
	logger  *wails.CustomLogger
}

// Settings handles settings values
type Settings struct {
	EnginePath              string
	BoardBackgroundColor    string
	BoardCoordinatesColor   string
	BoardLastMoveArrowColor string
	BoardWhiteCellsColor    string
	BoardBlackCellsColor    string
	BoardDndStartColor      string
	BoardDndEndColor        string
	BoardDndCrossColor      string
}

// WailsInit initialize the structure for Wails
func (settings *SettingsManager) WailsInit(r *wails.Runtime) error {
	settings.Runtime = r
	settings.logger = r.Log.New("AppSettings")
	return nil
}

// NewSettings builds a new settings instance
func NewSettings() *SettingsManager {
	return &SettingsManager{}
}

// Load returns a JSON string of the settings
// Tries to load from the configuration folder in home directory,
// or if it fails from the folder of program execution.
// If could not get settings from file, will use a default set of values.
// [output] string: the JSON string, or the error if the settings could not
// be converted into JSON
// #ConversionToJsonError : could not convert the settings to JSON
func (settings *SettingsManager) Load() string {

	defaultSettings := Settings{
		EnginePath:              "",
		BoardBackgroundColor:    "#06f",
		BoardCoordinatesColor:   "#f90",
		BoardLastMoveArrowColor: "#cde",
		BoardWhiteCellsColor:    "#fc0",
		BoardBlackCellsColor:    "#a33",
		BoardDndStartColor:      "#12FC23",
		BoardDndEndColor:        "#FC2312",
		BoardDndCrossColor:      "#898923",
	}

	dir, err := settings.Runtime.FileSystem.HomeDir()
	if err != nil {
		path, err := os.Executable()
		if err != nil {
			usedSettings, err := json.Marshal(defaultSettings)
			if err != nil {
				return "#ConversionToJsonError"
			}
			return string(usedSettings)
		}
		dir = path
	}

	filePath := path.Join(dir, "ChessExercisesOrganizerData", "config.txt")
	text, err := ioutil.ReadFile(filePath)
	if err != nil {
		usedSettings, err := json.Marshal(defaultSettings)
		if err != nil {
			return "#ConversionToJsonError"
		}
		return string(usedSettings)
	}

	usedSettings, err := json.Marshal(string(text))
	if err != nil {
		return "#ConversionToJsonError"
	}
	return string(usedSettings)
}

// Save attempts to save given update to the settings file.
// updateJSON JSON: a json object having the fields of Settings type.
// [output] string: the error if any.
// #ConversionError: failed to convert the given json
// #FileSavingError: failed to save content into file
func (settings *SettingsManager) Save(updateJSON map[string]interface{}) string {
	var result Settings

	err := mapstructure.Decode(updateJSON, &result)
	if err != nil {
		return "#ConversionError"
	}

	resultText, err := json.Marshal(result)
	if err != nil {
		return "#ConversionError"
	}

	var dir string

	dir, err = settings.Runtime.FileSystem.HomeDir()
	if err != nil {
		dir, err = os.Executable()
		if err != nil {
			return "#FileSavingError"
		}
	}

	settingsFolder := "ChessExercisesOrganizerData"
	err = os.MkdirAll(path.Join(dir, settingsFolder), os.ModePerm)
	if err != nil {
		return "#FileSavingError"
	}

	filePath := path.Join(dir, settingsFolder, "config.txt")
	file, err := os.Create(filePath)
	if err != nil {
		return "#FileSavingError"
	}
	defer file.Close()

	_, err = file.WriteString(string(resultText))
	if err != nil {
		return "#FileSavingError"
	}

	file.Sync()
	return ""
}
