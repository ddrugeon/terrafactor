/*
MIT License

# Copyright 2022 - © David Drugeon-Hamon

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
package state_test

import (
	"strings"
	"testing"

	"github.com/ddrugeon/terrafactor/internal/state"

	"github.com/stretchr/testify/assert"
)

func TestFromReader(t *testing.T) {
	t.Run("Empty state file", func(t *testing.T) {
		got := state.FromReader(strings.NewReader("{}"))
		wanted := state.TerraformState{}
		assert.Equal(t, got, wanted, "Terraform state read is not valid")
	})

	t.Run("Neither resource nor outputs in state file", func(t *testing.T) {
		got := state.FromReader(strings.NewReader("{\"terraform_version\": \"1.1.9\", \"version\": 4, \"lineage\": \"1681d92d-0964-f5eb-73d4-6e2dfa00baca\"}"))
		wanted := state.TerraformState{Version: 4, TerraformVersion: "1.1.9", Lineage: "1681d92d-0964-f5eb-73d4-6e2dfa00baca"}
		assert.Equal(t, got, wanted, "Terraform state read is not valid")
	})

	t.Run("No resource in state file", func(t *testing.T) {
		got := state.FromReader(strings.NewReader("{\n  \"version\": 4,\n  \"terraform_version\": \"1.1.9\",\n  \"serial\": 454,\n  \"lineage\": \"16\",\n  \"outputs\": {\n    \"datadog_synthetics_test_count\": {\n      \"value\": 387,\n      \"type\": \"number\"\n    }\n  },\n  \"resources\": []\n}\n"))
		wanted := state.TerraformState{Version: 4, Serial: 454, TerraformVersion: "1.1.9", Lineage: "16", Resources: []state.TerraformResource{}, Outputs: map[string]state.TerraformOutputValue{"datadog_synthetics_test_count": {Value: 387, Type: "number", Description: "", Sensitive: false}}}
		assert.EqualValues(t, got, wanted, "Terraform state read is not valid")
	})

	t.Run("resource and outputs in state file", func(t *testing.T) {
		got := state.FromReader(strings.NewReader("{\n  \"version\": 4,\n  \"terraform_version\": \"1.1.9\",\n  \"serial\": 454,\n  \"lineage\": \"16\",\n  \"outputs\": {\n    \"datadog_synthetics_test_count\": {\n      \"value\": 387,\n      \"type\": \"number\"\n    }\n  },\n  \"resources\": [\n    {\n      \"mode\": \"data\",\n      \"type\": \"aws_caller_identity\",\n      \"name\": \"current\",\n      \"provider\": \"provider[\\\"registry.terraform.io/hashicorp/aws\\\"]\",\n      \"instances\": [\n        {\n          \"schema_version\": 0,\n          \"attributes\": {\n            \"account_id\": \"123456789\",\n            \"arn\": \"arn:aws:sts::123456789:assumed-role/test/instance\",\n            \"id\": \"123456789\",\n            \"user_id\": \"ABC:instance\"\n          },\n          \"sensitive_attributes\": []\n        }\n      ]\n    },\n    {\n      \"module\": \"module.test\",\n      \"mode\": \"managed\",\n      \"type\": \"null_resource\",\n      \"name\": \"dummy_trigger\",\n      \"provider\": \"provider[\\\"registry.terraform.io/hashicorp/null\\\"]\",\n      \"instances\": [\n        {\n          \"schema_version\": 0,\n          \"attributes\": {\n            \"id\": \"123\"\n          },\n          \"sensitive_attributes\": [],\n          \"private\": \"AAA==\"\n        }\n      ]\n    }\n  ]\n}\n"))
		wanted := state.TerraformState{
			Version:          4,
			Serial:           454,
			TerraformVersion: "1.1.9",
			Lineage:          "16",
			Resources: []state.TerraformResource{
				{
					Mode:     "data",
					Type:     "aws_caller_identity",
					Name:     "current",
					Provider: "provider[\"registry.terraform.io/hashicorp/aws\"]",
					Instances: []state.TerraformResourceValue{
						{
							SchemaVersion: 0,
							Attributes: map[string]interface{}{
								"account_id": "123456789",
								"arn":        "arn:aws:sts::123456789:assumed-role/test/instance",
								"id":         "123456789",
								"user_id":    "ABC:instance",
							},
							SensitiveAttributes: []interface{}{},
						},
					},
				},
				{
					Module:   "module.test",
					Mode:     "managed",
					Type:     "null_resource",
					Name:     "dummy_trigger",
					Provider: "provider[\"registry.terraform.io/hashicorp/null\"]",
					Instances: []state.TerraformResourceValue{
						{
							SchemaVersion: 0,
							Attributes: map[string]interface{}{
								"id": "123",
							},
							SensitiveAttributes: []interface{}{},
						},
					},
				},
			},
			Outputs: map[string]state.TerraformOutputValue{
				"datadog_synthetics_test_count": {
					Value:       387,
					Type:        "number",
					Description: "",
					Sensitive:   false,
				},
			},
		}
		assert.EqualValues(t, got, wanted, "Terraform state read is not valid")
	})
}

func TestTerraformState_ListResources(t *testing.T) {
	t.Run("Should returns an empty resource list when empty state file is given", func(t *testing.T) {
		input := state.TerraformState{}
		expected := []state.TerraformResource{}

		assert.Equal(t, expected, input.ListResources(state.ResourceFilter{}))
	})
	t.Run("Should returns an empty resource list when no resources in state", func(t *testing.T) {
		input := state.TerraformState{
			Resources: []state.TerraformResource{},
		}
		expected := []state.TerraformResource{}

		assert.Equal(t, expected, input.ListResources(state.ResourceFilter{}))
	})
	t.Run("Should returns an empty resource list when no resources in state and filter is specified", func(t *testing.T) {
		input := state.TerraformState{
			Resources: []state.TerraformResource{},
		}
		filter := state.ResourceFilter{
			Type:     "null_resource",
			Module:   "module.test",
			Mode:     "managed",
			Name:     "dummy_trigger",
			Provider: "provider[\"registry.terraform.io/hashicorp/null\"]",
		}

		expected := []state.TerraformResource{}

		assert.Equal(t, expected, input.ListResources(filter))
	})
	t.Run("Should returns an empty resource list when non matching resources in state with specified filter", func(t *testing.T) {
		input := state.TerraformState{
			Version:          4,
			Serial:           454,
			TerraformVersion: "1.1.9",
			Lineage:          "16",
			Resources: []state.TerraformResource{
				{
					Mode:     "data",
					Type:     "aws_caller_identity",
					Name:     "current",
					Provider: "provider[\"registry.terraform.io/hashicorp/aws\"]",
					Instances: []state.TerraformResourceValue{
						{
							SchemaVersion: 0,
							Attributes: map[string]interface{}{
								"account_id": "123456789",
								"arn":        "arn:aws:sts::123456789:assumed-role/test/instance",
								"id":         "123456789",
								"user_id":    "ABC:instance",
							},
							SensitiveAttributes: []interface{}{},
						},
					},
				},
			},
			Outputs: map[string]state.TerraformOutputValue{
				"datadog_synthetics_test_count": {
					Value:       387,
					Type:        "number",
					Description: "",
					Sensitive:   false,
				},
			},
		}
		filter := state.ResourceFilter{
			Type:     "null_resource",
			Module:   "module.test",
			Mode:     "managed",
			Name:     "dummy_trigger",
			Provider: "provider[\"registry.terraform.io/hashicorp/null\"]",
		}

		expected := []state.TerraformResource{}

		assert.Equal(t, expected, input.ListResources(filter))
	})
	t.Run("Should returns a non empty resource list when no filter", func(t *testing.T) {
		input := state.TerraformState{
			Version:          4,
			Serial:           454,
			TerraformVersion: "1.1.9",
			Lineage:          "16",
			Resources: []state.TerraformResource{
				{
					Mode:     "data",
					Type:     "aws_caller_identity",
					Name:     "current",
					Provider: "provider[\"registry.terraform.io/hashicorp/aws\"]",
					Instances: []state.TerraformResourceValue{
						{
							SchemaVersion: 0,
							Attributes: map[string]interface{}{
								"account_id": "123456789",
								"arn":        "arn:aws:sts::123456789:assumed-role/test/instance",
								"id":         "123456789",
								"user_id":    "ABC:instance",
							},
							SensitiveAttributes: []interface{}{},
						},
					},
				},
				{
					Module:   "module.test",
					Mode:     "managed",
					Type:     "null_resource",
					Name:     "dummy_trigger",
					Provider: "provider[\"registry.terraform.io/hashicorp/null\"]",
					Instances: []state.TerraformResourceValue{
						{
							SchemaVersion: 0,
							Attributes: map[string]interface{}{
								"id": "123",
							},
							SensitiveAttributes: []interface{}{},
						},
					},
				},
			},
			Outputs: map[string]state.TerraformOutputValue{
				"datadog_synthetics_test_count": {
					Value:       387,
					Type:        "number",
					Description: "",
					Sensitive:   false,
				},
			},
		}
		filter := state.ResourceFilter{
			Type:     "null_resource",
			Module:   "module.test",
			Mode:     "managed",
			Name:     "dummy_trigger",
			Provider: "provider[\"registry.terraform.io/hashicorp/null\"]",
		}

		expected := []state.TerraformResource{
			{
				Module:   "module.test",
				Mode:     "managed",
				Type:     "null_resource",
				Name:     "dummy_trigger",
				Provider: "provider[\"registry.terraform.io/hashicorp/null\"]",
				Instances: []state.TerraformResourceValue{
					{
						SchemaVersion: 0,
						Attributes: map[string]interface{}{
							"id": "123",
						},
						SensitiveAttributes: []interface{}{},
					},
				},
			},
		}

		assert.Equal(t, expected, input.ListResources(filter))
	})

}
