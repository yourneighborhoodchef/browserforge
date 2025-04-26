package bayesian

import (
   "encoding/json"
   "fmt"

   "github.com/yourneighborhoodchef/browserforge/internal/data"
)

// networkDefinition represents the JSON structure of the Bayesian network.
type networkDefinition struct {
   Nodes []nodeDefinition `json:"nodes"`
}

// BayesianNetwork represents a Bayesian network.
type BayesianNetwork struct {
   nodesInOrder []*BayesianNode
   nodesByName  map[string]*BayesianNode
}

// LoadInputNetwork loads the input Bayesian network from embedded data.
func LoadInputNetwork() (*BayesianNetwork, error) {
   return loadNetwork(data.InputNetwork)
}

// LoadHeaderNetwork loads the header Bayesian network from embedded data.
func LoadHeaderNetwork() (*BayesianNetwork, error) {
   return loadNetwork(data.HeaderNetwork)
}

// LoadFingerprintNetwork loads the fingerprint Bayesian network from embedded data.
func LoadFingerprintNetwork() (*BayesianNetwork, error) {
   return loadNetwork(data.FingerprintNetwork)
}

// loadNetwork parses raw JSON bytes into a BayesianNetwork.
func loadNetwork(raw []byte) (*BayesianNetwork, error) {
   var def networkDefinition
   if err := json.Unmarshal(raw, &def); err != nil {
       return nil, fmt.Errorf("failed to unmarshal network: %w", err)
   }
   bn := &BayesianNetwork{
       nodesByName: make(map[string]*BayesianNode, len(def.Nodes)),
   }
   for _, nd := range def.Nodes {
       node := &BayesianNode{def: nd}
       bn.nodesInOrder = append(bn.nodesInOrder, node)
       bn.nodesByName[node.Name()] = node
   }
   return bn, nil
}

// GenerateSample randomly samples a full set of values from the network.
func (bn *BayesianNetwork) GenerateSample(inputValues map[string]string) (map[string]string, error) {
   sample := make(map[string]string)
   for k, v := range inputValues {
       sample[k] = v
   }
   for _, node := range bn.nodesInOrder {
       if _, exists := sample[node.Name()]; !exists {
           val, err := node.Sample(sample)
           if err != nil {
               return nil, err
           }
           sample[node.Name()] = val
       }
   }
   return sample, nil
}