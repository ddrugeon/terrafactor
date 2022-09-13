// Package resources create cli commands to manage ressources found in terraform state file
/*
MIT License

Copyright 2022 - © David Drugeon-Hamon

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package resources

import (
	"errors"
	"fmt"
	"os"

	"github.com/ddrugeon/terrafactor/cmd/options"
	"github.com/ddrugeon/terrafactor/internal/state"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

// NewListCommand creates a new `resources list` command
func NewListCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List all the boards",
		Run:   listResources,
		Args:  cobra.NoArgs,
	}

	command.PersistentFlags().StringVarP(&options.TerraformStateFilePath, options.ArgTFStateFile, options.Args[options.ArgTFStateFile].Short, options.Args[options.ArgTFStateFile].DefaultValue, options.Args[options.ArgTFStateFile].Description)
	err := command.MarkPersistentFlagRequired(options.ArgTFStateFile)
	if err != nil {
		return nil
	}
	command.PersistentFlags().StringVarP(&options.ResourceFilterString, options.ArgResourceFilter, options.Args[options.ArgResourceFilter].Short, options.Args[options.ArgResourceFilter].DefaultValue, options.Args[options.ArgResourceFilter].Description)

	return command
}

func listResources(cmd *cobra.Command, args []string) {
	filter, err := state.CreateResourceFilterFromString(options.ResourceFilterString)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}

	if _, err := os.Stat(options.TerraformStateFilePath); errors.Is(err, os.ErrNotExist) {
		pterm.Error.Printf("%s not found\n", options.TerraformStateFilePath)
		os.Exit(1)
	}

	file, err := os.Open(options.TerraformStateFilePath)
	if err != nil {
		pterm.Error.Printf("Error opening terraform state file %s - %s\n", options.TerraformStateFilePath, err)
		os.Exit(1)
	}
	defer file.Close()

	terraformState := state.FromReader(file)
	terraformResources := terraformState.ListResources(*filter)

	panels := pterm.Panels{
		{{Data: pterm.Yellow("\nTerraform Version:\nProcessed State file:")}, {Data: fmt.Sprintf("\n%s\n%s", terraformState.TerraformVersion, options.TerraformStateFilePath)}},
	}
	_ = pterm.DefaultPanel.WithPanels(panels).WithPadding(5).Render()

	resourcesList := pterm.LeveledList{pterm.LeveledListItem{Level: 0, Text: pterm.Yellow("Resources")}}
	for _, resource := range terraformResources {
		if resource.Mode == "managed" {
			resourceType := resource.Type
			resourceName := resource.Name
			for _, instance := range resource.Instances {
				indexKey := instance.IndexKey
				resourceLocation := fmt.Sprintf("%s.%s", resourceType, resourceName)
				if len(resource.Module) > 0 {
					resourceLocation = fmt.Sprintf("%s.%s", resource.Module, resourceLocation)
				}
				if len(indexKey) > 0 {
					resourceLocation = fmt.Sprintf("%s[\"%s\"]", resourceLocation, indexKey)
				}

				resourcesList = append(resourcesList, pterm.LeveledListItem{Level: 1, Text: resourceLocation})

			}
		}
	}

	root := putils.TreeFromLeveledList(resourcesList)
	_ = pterm.DefaultTree.WithRoot(root).Render()
}
