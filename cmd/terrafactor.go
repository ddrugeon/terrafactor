// Package cmd list cli commands
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
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ddrugeon/terrafactor/cmd/options"

	"github.com/ddrugeon/terrafactor/cmd/resources"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

const (
	envPrefix = "TFCT"
)

// NewRootCommand builds the main cli application and
// adds children command hierarchy
func NewRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   options.ApplicationName,
		Short: options.ApplicationShort,
		Long:  options.ApplicationLong,
	}

	command.AddCommand(
		resources.NewResourceCommand(),
		NewVersionCommand(),
	)

	return command
}

// Execute is the main entry point of this command
func Execute() {
	initConfig()

	cmd := NewRootCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func init() {
	cobra.OnInitialize(initConfig)
}
