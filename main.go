package main

import (
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/sebastianxyzsss/tardigrade/globals"
	"github.com/sebastianxyzsss/tardigrade/logger"
	"github.com/sebastianxyzsss/tardigrade/parser"
	"github.com/sebastianxyzsss/tardigrade/reader"
	"github.com/sebastianxyzsss/tardigrade/selecter"
	"github.com/spf13/viper"
)

type Cli struct {
	FlatParse bool     `short:"f" help:"Flat parse mode, instead of hierarchical parse"`
	New       bool     `short:"n" help:"Add new user home content file"`
	Init      bool     `short:"i" help:"Add new content file in local directory"`
	Tags      bool     `short:"t" help:"Make strings in command, Tag filters, default"`
	All       bool     `short:"a" help:"Make strings in command, filter Anything, default"`
	Files     bool     `short:"s" help:"Make strings in command, files to read, default"`
	Paths     []string `arg:"" optional:"" name:"path" help:"Paths to remove." type:"path"`
}

var CLI Cli

var ll = logger.SetupLog()

func main() {

	ctx := kong.Parse(&CLI)

	switch ctx.Command() {
	case "<path>":
		getOptions(&CLI)
		runTardigrade(CLI.Paths)
	default:
		getOptions(&CLI)
		runTardigrade(nil)
	}
}

func getOptions(object interface{}) {
	cli, ok := object.(*Cli)
	if ok {
		if cli.FlatParse {
			ll.Debug().Msg("flat parse mode enabled")
			globals.FlatParse = true
		}
		if cli.Files {
			ll.Debug().Msg("set strings as files to read")
			globals.FilterAction = globals.FilterFiles
		}
		if cli.Tags {
			ll.Debug().Msg("set strings as Tag filters, set to parse flat")
			globals.FlatParse = true
			globals.FilterAction = globals.FilterTags
			for _, path := range CLI.Paths {
				globals.FilterStrings = append(globals.FilterStrings, filepath.Base(path))
			}
		}
		if cli.All {
			ll.Debug().Msg("set strings as All filters, in anything, set to parse flat")
			globals.FlatParse = true
			globals.FilterAction = globals.FilterAnything
			for _, path := range CLI.Paths {
				globals.FilterStrings = append(globals.FilterStrings, filepath.Base(path))
			}
		}
		if cli.Init {
			ll.Info().Msg("creating new local file")
			reader.CreateNewLocalContentFile()
			os.Exit(0)
		}
		if cli.New {
			ll.Info().Msg("creating new user home file")
			reader.CreateNewUserContentFile()
			os.Exit(0)
		}
	} else {
		ll.Error().Msg("Invalid cli type")
	}
}

type conf struct {
	Settings globals.Settings
}

func createSettings() (*globals.Settings, error) {
	ll.Debug().Msg("createSettings")

	userHomeDirName, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := userHomeDirName + "/.tardigrade/tardisettings.yml"

	viper.SetConfigFile(configPath)

	if reader.DoesFileNotExist(configPath) {
		ll.Debug().Msg("seeting file does not exist, creating")

		viper.SetDefault("settings", globals.Settings{
			Height:           11,
			IndicatorStyle:   "80",
			FooterKeyMaxSize: 16,
			HistorySize:      11,
			LogLevel:         "info",
		})

		viper.WriteConfig()
	}

	c := &conf{}

	ll.Debug().Msg("reading " + configPath)

	viper.ReadInConfig()
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}

	return &c.Settings, nil
}

func runTardigrade(strsToRead []string) {
	ll.Debug().Msg("-------------------------------------------------------- tardigrade ..")

	settings, err := createSettings()
	if err != nil {
		ll.Error().Err(err).Msg("there is an error in the settings")
		os.Exit(1)
	}

	globals.ChildKeyMaxSize = settings.FooterKeyMaxSize
	globals.HistorySize = settings.HistorySize

	if settings.LogLevel == "debug" {
		logger.SetLogLevelDebug()
	} else {
		logger.SetLogLevelInfo()
	}

	yamlAsMap := reader.GetRawMapContent(strsToRead)

	rootElement := parser.NewElement("root", false, nil)
	err = parser.MainRecurseMap(*yamlAsMap, rootElement)
	if err != nil {
		ll.Error().Err(err).Msg("No children were found after filtering")
		os.Exit(1)
	}
	ll.Debug().Msg("-----------")
	ll.Debug().Msg("-----------")

	selecter.Chooser(rootElement, settings)

	ll.Debug().Msg("------------------------------------------------------")
}
