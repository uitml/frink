package types

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// StringArray is a proxy type to support umarshalling YAML nodes that are scalars or sequences of type string.
type StringArray []string

// UnmarshalYAML implements behavior to handle nodes that are scalars or sequences of type string.
// TODO: Refactor this; don't like the way it looks and "feels"...
func (array *StringArray) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.SequenceNode:
		var sequence []string
		if err := node.Decode(&sequence); err != nil {
			return err
		}
		*array = sequence
	case yaml.ScalarNode:
		var scalar string
		if err := node.Decode(&scalar); err != nil {
			return err
		}
		*array = []string{scalar}
	default:
		return fmt.Errorf("kind must be scalar or sequence")
	}

	return nil
}
