package plugins

import (
	"embed"
	"fmt"
	"os"
	"path"

	"github.com/jazzathedev/bam/internal/plugin"
	"github.com/jazzathedev/bam/internal/setup"
)

//go:embed builtin/*.toml
var BuiltinPlugins embed.FS

func loadPlugins(pluginStrings []string) ([]plugin.PluginConfig, error) {
	var pluginsStruct []plugin.PluginConfig

	for _, pluginString := range pluginStrings {

		pluginStruct, err := plugin.LoadPlugin(pluginString)
		if err != nil {
			return nil, err
		}

		pluginsStruct = append(pluginsStruct, pluginStruct)
	}

	return pluginsStruct, nil
}

func LoadBuiltinPlugins() ([]plugin.PluginConfig, error) {
	dirEntries, err := BuiltinPlugins.ReadDir("builtin")
	if err != nil {
		return nil, err
	}

	var pluginStrings []string

	for _, dirEntry := range dirEntries {
		pluginString, err := BuiltinPlugins.ReadFile(path.Join("builtin", dirEntry.Name()))
		if err != nil {
			return nil, fmt.Errorf("Unable to load builtin plugins: %w", err)
		}

		pluginStrings = append(pluginStrings, string(pluginString))
	}

	return loadPlugins(pluginStrings)

}

func LoadUserPlugins() ([]plugin.PluginConfig, error) {
	bam, err := setup.BamDir()
	if err != nil {
		return nil, err
	}

	var pluginStrings []string

	dirEntries, err := os.ReadDir(path.Join(bam, "plugins/user"))
	if err != nil {
		return nil, fmt.Errorf("Unable to load user plugins folder: %w", err)
	}

	for _, dirEntry := range dirEntries {
		pluginPath := path.Join(bam, "plugins/user", dirEntry.Name())

		pluginString, err := os.ReadFile(pluginPath)
		if err != nil {
			return nil, fmt.Errorf("Unable to load user plugin toml file: %w", err)
		}

		pluginStrings = append(pluginStrings, string(pluginString))
	}

	return loadPlugins(pluginStrings)
}
