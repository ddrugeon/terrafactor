// Package options gives the list of argument options for terrafactor cli
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
package options

type arguments struct {
	Description  string
	Short        string
	DefaultValue string
}

const (
	// ApplicationName is the name of this CLI
	ApplicationName = "terrafactor"
	// ApplicationShort is the short description of this CLI
	ApplicationShort = "terrafactor CLI"
	// ApplicationLong is the long description of this CLI
	ApplicationLong = `terrafactor generates terraform move statement from current terraform state file.`

	// ArgTFStateFile represents the name of option to specify terraform state file
	ArgTFStateFile = "tfstate"

	// ArgResourceFilter is the name of flag to specify a resource string
	ArgResourceFilter = "filter"
)

// Args represents lists different options for one argument (Description, Short, DefaultValue)
var Args = map[string]arguments{
	ArgTFStateFile: {
		Description:  "(required) Path of the terraform state in json",
		Short:        "t",
		DefaultValue: "",
	},
	ArgResourceFilter: {
		Description:  "(optional) Filter string to apply - Example: resource.datadog_synthetics_private_location.main",
		Short:        "f",
		DefaultValue: "",
	},
}

// TerraformStateFilePath is the path where terraform state file can be found
var TerraformStateFilePath string

// ResourceFilterString is a string representation for a filter
var ResourceFilterString string
