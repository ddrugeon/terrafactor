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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

// TerraformOutputValue represents a value of terraform output.
type TerraformOutputValue struct {
	Sensitive   bool   `json:"sensitive"`
	Type        string `json:"type"`
	Value       int    `json:"value"`
	Description string `json:"description"`
}

// TerraformResourceValue represents a value of terraform resource (or module).
type TerraformResourceValue struct {
	IndexKey            string                 `json:"index_key"`
	SchemaVersion       int                    `json:"schema_version"`
	Attributes          map[string]interface{} `json:"attributes"`
	SensitiveAttributes []interface{}          `json:"sensitive_attributes"`
}

// TerraformResource represents terraform resource (or module).
type TerraformResource struct {
	Module    string                   `json:"module"`
	Mode      string                   `json:"mode"`
	Type      string                   `json:"type"`
	Name      string                   `json:"name"`
	Provider  string                   `json:"provider"`
	Instances []TerraformResourceValue `json:"instances"`
}

// TerraformState represents terraform state.
type TerraformState struct {
	Version          int                             `json:"version"`
	TerraformVersion string                          `json:"terraform_version"`
	Serial           int                             `json:"serial"`
	Lineage          string                          `json:"lineage"`
	Outputs          map[string]TerraformOutputValue `json:"outputs"`
	Resources        []TerraformResource             `json:"resources"`
}

func (resource TerraformResource) String() string {
	prefix := ""
	if strings.TrimSpace(resource.Module) != "" {
		prefix = fmt.Sprintf("%s.", resource.Module)
	}

	return fmt.Sprintf("%s%s.%s", prefix, resource.Type, resource.Name)
}

// FromReader unmarshall terraform state from a reader.
func FromReader(reader io.Reader) TerraformState {
	terraformState := TerraformState{}
	decoder := json.NewDecoder(reader)

	for {
		var err = decoder.Decode(&terraformState)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error decoding json data:", err)
			break
		}
	}

	return terraformState
}

// ListResources returns a map with two entries: Llist of resources. Resources and Modules can be
// filtered with ResourceFilter
func (s TerraformState) ListResources(filter ResourceFilter) []TerraformResource {
	output := []TerraformResource{}
	for _, resource := range s.Resources {
		if filter.Matches(resource) {
			output = append(output, resource)
		}
	}
	return output
}

// GenerateMovedStatement generates terraform moved statement for a resource to a newLocation
func GenerateMovedStatement(resource TerraformResource, newLocation string) string {
	if len(resource.Instances) == 0 {
		return ""
	}
	output := ""
	for _, instance := range resource.Instances {
		if strings.TrimSpace(instance.IndexKey) == "" {
			output += fmt.Sprintf("moved {\n  from = %s\n  to   = %s\n}\n\n", resource, newLocation)
		} else {
			output += fmt.Sprintf("moved {\n  from = %s[\"%s\"]\n  to   = %s[\"%s\"]\n}\n\n", resource, instance.IndexKey, newLocation, instance.IndexKey)
		}
	}

	return output
}
