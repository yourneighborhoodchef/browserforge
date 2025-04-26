package bayesian

import (
   "errors"
   "fmt"
   "math/rand"
)

// nodeDefinition represents the JSON structure of a Bayesian node.
type nodeDefinition struct {
   Name                     string                 `json:"name"`
   ParentNames              []string               `json:"parentNames"`
   ConditionalProbabilities map[string]interface{} `json:"conditionalProbabilities"`
}

// BayesianNode represents a node in a Bayesian network.
type BayesianNode struct {
   def nodeDefinition
}

// Name returns the name of the node.
func (n *BayesianNode) Name() string {
   return n.def.Name
}

// ParentNames returns the names of parent nodes.
func (n *BayesianNode) ParentNames() []string {
   return n.def.ParentNames
}

// Sample randomly samples a value for this node given parent values.
func (n *BayesianNode) Sample(parentValues map[string]string) (string, error) {
   probs, err := n.probabilitiesGiven(parentValues)
   if err != nil {
       return "", err
   }
   // Prepare for weighted sampling.
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
   // Fallback to first value.
   return values[0], nil
}

// probabilitiesGiven computes the conditional probabilities for this node given parentValues.
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
       // Otherwise, keep current.
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