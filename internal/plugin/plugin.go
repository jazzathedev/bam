package plugin

import (
	"github.com/BurntSushi/toml"
)

type PluginConfig struct {
	Name        string   `toml:"name"`
	Aliases     []string `toml:"aliases"`
	Description string   `toml:"description"`
	Schema      int      `toml:"schema"`
	Versions    Versions `toml:"versions"`
	Download    Download `toml:"download"`
	Platform    Platform `toml:"platform"`
	Install     Install  `toml:"install"`
}

type Versions struct {
	ListURL     string `toml:"list_url"`
	ListPath    string `toml:"list_path"`
	StripPrefix string `toml:"strip_prefix"`
}

type Download struct {
	URL        string `toml:"url"`
	HashURL    string `toml:"hash_url"`
	HashAlgo   string `toml:"hash_algo"`
	HashFormat string `toml:"hash_format"`
}

type Platform struct {
	OSMap   map[string]string `toml:"os_map"`
	ArchMap map[string]string `toml:"arch_map"`
	ExtMap  map[string]string `toml:"ext_map"`
}

type Install struct {
	StripComponents bool  `toml:"strip_components"`
	Bin             []Bin `toml:"bin"`
}

type Bin struct {
	Name string              `toml:"name"`
	Run  map[string][]string `toml:"run"`
}

func LoadPlugin(pluginFileString string) (PluginConfig, error) {
	var plugin PluginConfig

	_, err := toml.Decode(pluginFileString, &plugin)
	if err != nil {
		return PluginConfig{}, err
	}

	return plugin, nil
}
