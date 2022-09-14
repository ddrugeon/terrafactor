// Package state list all related resource to terraform state file
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
package state

import (
	"errors"
	"fmt"
	"strings"
)

const filterErrorMsg = "Filter must conform pattern: type.name or module.module_name.type.name"

// ResourceFilter is a struct giving criteria to filter resources
type ResourceFilter struct {
	Type     string
	Mode     string
	Name     string
	Provider string
	Module   string
}

func (f ResourceFilter) matchType(value string) bool {
	return (strings.TrimSpace(f.Type) == "") || (value == f.Type)
}

func (f ResourceFilter) matchMode(value string) bool {
	return (strings.TrimSpace(f.Mode) == "") || (value == f.Mode)
}

func (f ResourceFilter) matchModule(value string) bool {
	return (strings.TrimSpace(f.Module) == "") || (value == f.Module)
}

func (f ResourceFilter) matchName(value string) bool {
	return (strings.TrimSpace(f.Name) == "") || (value == f.Name)
}

func (f ResourceFilter) matchProvider(value string) bool {
	return (strings.TrimSpace(f.Provider) == "") || (value == f.Provider)
}

// Matches return true if all ResourceFilter properties matches specified resource
func (f ResourceFilter) Matches(resource TerraformResource) bool {
	if f.isEmptyFilter() {
		return true
	}

	return f.matchType(resource.Type) &&
		f.matchMode(resource.Mode) &&
		f.matchModule(resource.Module) &&
		f.matchProvider(resource.Provider) &&
		f.matchName(resource.Name)
}

func (f ResourceFilter) isEmptyFilter() bool {
	return f.Mode == "" && f.Type == "" && f.Provider == "" && f.Module == "" && f.Name == ""
}

// CreateResourceFilterFromString create a ResourceFilter from a string
func CreateResourceFilterFromString(filter string) (*ResourceFilter, error) {
	output := ResourceFilter{}
	if strings.TrimSpace(filter) != "" {
		fields := strings.Split(filter, ".")

		if strings.HasPrefix(filter, "module.") {
			if len(fields) != 4 {
				return nil, errors.New(filterErrorMsg)
			}
			output.Module = fmt.Sprintf("%s.%s", fields[0], fields[1])
			fields = fields[2:]
		} else if len(fields) != 2 {
			return nil, errors.New(filterErrorMsg)
		}

		for index, field := range fields {
			if strings.TrimSpace(field) == "" {
				return nil, errors.New(filterErrorMsg)
			}
			switch index {
			case 0:
				output.Type = field
			case 1:
				output.Name = field
			}
		}

	}

	return &output, nil
}
