package cli

import (
	"encoding/json"

	"github.com/kuttiproject/workspace"
)

var (
	data           *settingdata
	settingmanager workspace.Configmanager
)

type settingdata struct {
	settings map[string]string
}

func (dc *settingdata) Serialize() ([]byte, error) {
	return json.Marshal(dc.settings)
}

func (dc *settingdata) Deserialize(data []byte) error {
	loadeddefaults := make(map[string]string)
	err := json.Unmarshal(data, &loadeddefaults)
	if err == nil {
		dc.settings = loadeddefaults
	}

	return err
}

func (dc *settingdata) Setdefaults() {
	dc.settings = map[string]string{}
}

// Setting gets the value for the specified setting.
// If a value does not exist, an empty string is returned.
func Setting(configname string) (string, bool) {
	result, ok := data.settings[configname]
	return result, ok
}

// SetSetting sets the specified setting to the specified value.
func SetSetting(name string, value string) error {
	data.settings[name] = value
	return settingmanager.Save()
}

// RemoveSetting deletes the specified setting.
// If the setting does not exist, nothing happens.
func RemoveSetting(name string) error {
	delete(data.settings, name)
	return settingmanager.Save()
}

func settingnamefordefault(name string) string {
	return "default-" + name
}

// Default gets the value of a setting called default-<name>.
func Default(name string) (string, bool) {
	return Setting(settingnamefordefault(name))
}

// SetDefault sets a setting called default-<name> to the specified value.
func SetDefault(name string, value string) error {
	return SetSetting(settingnamefordefault(name), value)
}

// RemoveDefault deletes a setting called default-<name>.
// If the setting does not exist, nothing happens.
func RemoveDefault(name string) error {
	return RemoveSetting(settingnamefordefault(name))
}

// Settings returns the settings map.
func Settings() map[string]string {
	return data.settings
}

func init() {
	data = &settingdata{
		settings: map[string]string{},
	}

	var err error
	settingmanager, err = workspace.NewFileConfigmanager("config", data)

	// If config manager cannot be initialized, big trouble.
	if err != nil {
		panic(err)
	}
}
