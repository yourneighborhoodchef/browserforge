package bayesian

import (
	"errors"
	"fmt"
	"math/rand"
)

type nodeDefinition struct {
	Name                     string                 `json:"name"`
	ParentNames              []string               `json:"parentNames"`
	ConditionalProbabilities map[string]interface{} `json:"conditionalProbabilities"`
}

type BayesianNode struct {
	def nodeDefinition
}

func (n *BayesianNode) Name() string {
	return n.def.Name
}

func (n *BayesianNode) ParentNames() []string {
	return n.def.ParentNames
}

func (n *BayesianNode) Sample(parentValues map[string]string) (string, error) {
	probs, err := n.probabilitiesGiven(parentValues)
	if err != nil {
		return "", err
	}

	var values []string
	var weights []float64
	for v, w := range probs {
		values = append(values, v)
		weights = append(weights, w)
	}
	total := 0.0
	for _, w := range weights {
		total += w
	}
	if total <= 0 {
		return "", errors.New("total probability is zero")
	}
	target := rand.Float64() * total
	cum := 0.0
	for i, w := range weights {
		cum += w
		if target <= cum {
			return values[i], nil
		}
	}

	return values[0], nil
}

func (n *BayesianNode) probabilitiesGiven(parentValues map[string]string) (map[string]float64, error) {
	var current interface{} = n.def.ConditionalProbabilities
	for _, parent := range n.def.ParentNames {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type for conditional probabilities")
		}
		if val, exists := parentValues[parent]; exists {
			if deeper, ok := m["deeper"].(map[string]interface{}); ok {
				if next, found := deeper[val]; found {
					current = next
					continue
				}
			}
		}
		if skip, ok := m["skip"].(map[string]interface{}); ok {
			current = skip
		}

	}
	m, ok := current.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid final probabilities structure")
	}
	result := make(map[string]float64, len(m))
	for k, v := range m {
		switch x := v.(type) {
		case float64:
			result[k] = x
		case int:
			result[k] = float64(x)
		default:
			return nil, fmt.Errorf("unsupported probability type for %s: %T", k, v)
		}
	}
	return result, nil
}
