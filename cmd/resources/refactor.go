// Package resources list cli commands to list all ressources and modules found in terraform state file
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

	"github.com/spf13/cobra"
	"gitlab.com/david.drugeon-hamon/terrafactor/cmd/options"
	"gitlab.com/david.drugeon-hamon/terrafactor/internal/state"
)

var (
	oldLocation string
	newLocation string
)

// NewRefactorCommand is the command to generate terraform moved directives
func NewRefactorCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "refactor [flags] old_location new_location",
		Short: "Generate terraform moved directives",
		Run:   refactor,
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("Required arguments old_location or new_location are missing")
			}

			oldLocation = args[0]
			newLocation = args[1]
			return nil
		},
	}

	command.PersistentFlags().StringVarP(&options.TerraformStateFilePath, options.ArgTFStateFile, options.Args[options.ArgTFStateFile].Short, options.Args[options.ArgTFStateFile].DefaultValue, options.Args[options.ArgTFStateFile].Description)
	err := command.MarkPersistentFlagRequired(options.ArgTFStateFile)
	if err != nil {
		return nil
	}

	return command
}

func refactor(cmd *cobra.Command, args []string) {
	filter, err := state.CreateResourceFilterFromString(oldLocation)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := os.Stat(options.TerraformStateFilePath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s not found\n", options.TerraformStateFilePath)
		os.Exit(1)
	}

	file, err := os.Open(options.TerraformStateFilePath)
	if err != nil {
		fmt.Printf("Error opening terraform state file %s - %s\n", options.TerraformStateFilePath, err)
		os.Exit(1)
	}
	defer file.Close()

	terraformState := state.FromReader(file)
	terraformResources := terraformState.ListResources(*filter)

	for _, resource := range terraformResources {
		fmt.Println(state.GenerateMovedStatement(resource, newLocation))
	}
}
