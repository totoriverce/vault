package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*PluginListCommand)(nil)
var _ cli.CommandAutocomplete = (*PluginListCommand)(nil)

type PluginListCommand struct {
	*BaseCommand
}

func (c *PluginListCommand) Synopsis() string {
	return "Lists available plugins"
}

func (c *PluginListCommand) Help() string {
	helpText := `
Usage: vault plugin list [options] TYPE

  Lists available plugins registered in the catalog. This does not list whether
  plugins are in use, but rather just their availability. The last argument of
  type takes "auth", "database", or "secret".

  List all available plugins in the catalog:

      $ vault plugin list

  List all available database plugins in the catalog:

      $ vault plugin list database

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PluginListCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *PluginListCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *PluginListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PluginListCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	pluginTypeStr := strings.TrimSpace(args[0])
	pluginType := consts.PluginTypeUnknown
	if pluginTypeStr != "" {
		var err error
		pluginType, err = consts.ParsePluginType(pluginTypeStr)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing type: %s", err))
			return 2
		}
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	resp, err := client.Sys().ListPlugins(&api.ListPluginsInput{
		Type: pluginType,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing available plugins: %s", err))
		return 2
	}

	var names []string
	namesAdded := make(map[string]bool)
	for _, names := range resp.NamesByType {
		for _, name := range names {
			if ok := namesAdded[name]; !ok {
				names = append(names, name)
				namesAdded[name] = true
			}
		}
	}
	sort.Strings(names)

	switch Format(c.UI) {
	case "table":
		list := append([]string{"Plugins"}, names...)
		c.UI.Output(tableOutput(list, nil))
		return 0
	default:
		return OutputData(c.UI, names)
	}
}
