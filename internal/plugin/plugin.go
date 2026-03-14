package plugin

import (
	"github.com/BurntSushi/toml"
)

type PluginConfig struct {
	Name        string   `toml:"name"`
	Aliases     []string `toml:"aliases"`
	Description string   `toml:"description"`
	Versions    Versions `toml:"versions"`
	Download    Download `toml:"download"`
	Platform    Platform `toml:"platform"`
	Install     Install  `toml:"install"`
}

type Versions struct {
	ListURL     string   `toml:"list_url"`
	ListPath    string   `toml:"list_path"`
	StripPrefix string   `toml:"strip_prefix"`
	Channels    Channels `toml:"channels"`
}

type Channels struct {
	Latest string `toml:"latest"`
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
	StripComponents int      `toml:"strip_components"`
	Bin             []string `toml:"bin"`
}

func LoadPlugin(pluginFileString string) (PluginConfig, error) {
	var plugin PluginConfig

	_, err := toml.Decode(pluginFileString, &plugin)
	if err != nil {
		return PluginConfig{}, err
	}

	return plugin, nil
}
