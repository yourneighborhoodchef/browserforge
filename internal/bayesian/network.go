package bayesian

import (
	"encoding/json"
	"fmt"

	"github.com/yourneighborhoodchef/browserforge/internal/data"
)

type networkDefinition struct {
	Nodes []nodeDefinition `json:"nodes"`
}

type BayesianNetwork struct {
	nodesInOrder []*BayesianNode
	nodesByName  map[string]*BayesianNode
}

func LoadInputNetwork() (*BayesianNetwork, error) {
	return loadNetwork(data.InputNetwork)
}

func LoadHeaderNetwork() (*BayesianNetwork, error) {
	return loadNetwork(data.HeaderNetwork)
}

func LoadFingerprintNetwork() (*BayesianNetwork, error) {
	return loadNetwork(data.FingerprintNetwork)
}

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
