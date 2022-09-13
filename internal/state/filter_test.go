package state_test

import (
	"testing"

	"github.com/ddrugeon/terrafactor/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestThatFilterMatchReturnsTrueWhenResourceMatches(t *testing.T) {
	resource := state.TerraformResource{
		Type:     "null_resource",
		Module:   "module.test",
		Mode:     "managed",
		Name:     "dummy_trigger",
		Provider: "provider[\"registry.terraform.io/hashicorp/null\"]",
	}

	t.Run("Return True if Type matches", func(t *testing.T) {
		filter := state.ResourceFilter{
			Type: "null_resource",
		}
		assert.True(t, filter.Matches(resource))
	})

	t.Run("Return True if Mode matches", func(t *testing.T) {
		filter := state.ResourceFilter{
			Mode: "managed",
		}
		assert.True(t, filter.Matches(resource))
	})

	t.Run("Return True if Module matches", func(t *testing.T) {
		filter := state.ResourceFilter{
			Module: "module.test",
		}
		assert.True(t, filter.Matches(resource))
	})

	t.Run("Return True if Name matches", func(t *testing.T) {
		filter := state.ResourceFilter{
			Name: "dummy_trigger",
		}
		assert.True(t, filter.Matches(resource))
	})

	t.Run("Return True if Provider matches", func(t *testing.T) {
		filter := state.ResourceFilter{
			Provider: "provider[\"registry.terraform.io/hashicorp/null\"]",
		}
		assert.True(t, filter.Matches(resource))
	})

	t.Run("Filter on all resource properties", func(t *testing.T) {
		filter := state.ResourceFilter{
			Type:     "null_resource",
			Module:   "module.test",
			Mode:     "managed",
			Name:     "dummy_trigger",
			Provider: "provider[\"registry.terraform.io/hashicorp/null\"]",
		}
		assert.True(t, filter.Matches(resource))
	})

	t.Run("Should returns true if filter criterias are empty", func(t *testing.T) {
		filter := state.ResourceFilter{}
		assert.True(t, filter.Matches(resource))
	})
}
func TestThatFilterMatchReturnsFalseWhenResourceNotMatches(t *testing.T) {
	resource := state.TerraformResource{
		Type:     "aws_caller_identity",
		Module:   "",
		Mode:     "data",
		Name:     "dummy_trigger",
		Provider: "provider[\"registry.terraform.io/hashicorp/aws\"]",
	}

	t.Run("Filter only on resource type", func(t *testing.T) {
		filter := state.ResourceFilter{
			Type: "null_resource",
		}
		assert.False(t, filter.Matches(resource))
	})
	t.Run("Filter only on resource module name", func(t *testing.T) {
		filter := state.ResourceFilter{
			Type:   "aws_caller_identity",
			Module: "module",
		}
		assert.False(t, filter.Matches(resource))
	})
	t.Run("Filter only on resource mode", func(t *testing.T) {
		filter := state.ResourceFilter{
			Type:   "aws_caller_identity",
			Module: "",
			Mode:   "managed",
		}
		assert.False(t, filter.Matches(resource))
	})
	t.Run("Filter only on resource name", func(t *testing.T) {
		filter := state.ResourceFilter{
			Type:     "aws_caller_identity",
			Module:   "",
			Mode:     "data",
			Name:     "toto",
			Provider: "provider[\"registry.terraform.io/hashicorp/aws\"]",
		}
		assert.False(t, filter.Matches(resource))
	})
	t.Run("Filter only on resource provider name", func(t *testing.T) {
		filter := state.ResourceFilter{
			Type:     "aws_caller_identity",
			Module:   "",
			Mode:     "data",
			Name:     "dummy_trigger",
			Provider: "provider[\"registry.terraform.io/hashicorp/null\"]",
		}
		assert.False(t, filter.Matches(resource))
	})
}
func TestCreateResourceFilterFromString(t *testing.T) {
	t.Run("Empty string should create an empty ResourceFilter", func(t *testing.T) {
		input := ""
		output, err := state.CreateResourceFilterFromString(input)
		assert.Nil(t, err, "An error occured when reading filter from string")
		assert.Empty(t, output)
	})

	t.Run("If filter string pattern is not type.name, should returns an error", func(t *testing.T) {
		input := "type"
		output, err := state.CreateResourceFilterFromString(input)
		assert.NotNil(t, err, "Error should be not nil")
		assert.Nil(t, output, "Resource filter should be nil if pattern is not type.name.")

		input = "."
		output, err = state.CreateResourceFilterFromString(input)
		assert.NotNil(t, err, "Error should be not nil")
		assert.Nil(t, output, "Resource filter should be nil if pattern is not type.name.")

		input = "type."
		output, err = state.CreateResourceFilterFromString(input)
		assert.NotNil(t, err, "Error should be not nil")
		assert.Nil(t, output, "Resource filter should be nil if pattern is not type.name.")

		input = ".name"
		output, err = state.CreateResourceFilterFromString(input)
		assert.NotNil(t, err, "Error should be not nil")
		assert.Nil(t, output, "Resource filter should be nil if pattern is not type.name.")

		input = "type.name."
		output, err = state.CreateResourceFilterFromString(input)
		assert.NotNil(t, err, "Error should be not nil")
		assert.Nil(t, output, "Resource filter should be nil if pattern is not type.name.")
	})

	t.Run("If filter string pattern is not module.module_name.type.name, should returns an error", func(t *testing.T) {
		input := "module.name."
		output, err := state.CreateResourceFilterFromString(input)
		assert.NotNil(t, err, "Error should be not nil")
		assert.Nil(t, output, "Resource filter should be nil if pattern is not module.module_name.type.name.")

		input = "module.name.type"
		output, err = state.CreateResourceFilterFromString(input)
		assert.NotNil(t, err, "Error should be not nil")
		assert.Nil(t, output, "Resource filter should be nil if pattern is not module.module_name.type.name.")

		input = "module.name.type."
		output, err = state.CreateResourceFilterFromString(input)
		assert.NotNil(t, err, "Error should be not nil")
		assert.Nil(t, output, "Resource filter should be nil if pattern is not module.module_name.type.name.")

		input = "module.name.type.."
		output, err = state.CreateResourceFilterFromString(input)
		assert.NotNil(t, err, "Error should be not nil")
		assert.Nil(t, output, "Resource filter should be nil if pattern is not module.module_name.type.name.")
	})

	t.Run("If filter string pattern is correct, should returns a valid ResourceFilter", func(t *testing.T) {
		input := "type.name"
		expected := state.ResourceFilter{
			Type: "type",
			Name: "name",
		}
		output, err := state.CreateResourceFilterFromString(input)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, *output, expected, "Resource filter should be not nil if pattern is type.name or module.module_name.type.name")

		input = "module.module_name.type.name"
		expected = state.ResourceFilter{
			Module: "module.module_name",
			Type:   "type",
			Name:   "name",
		}
		output, err = state.CreateResourceFilterFromString(input)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, expected, *output, "Resource filter should be not nil if pattern is type.name or module.module_name.type.name")
	})
}
