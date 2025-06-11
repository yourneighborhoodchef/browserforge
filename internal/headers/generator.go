package headers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourneighborhoodchef/browserforge/internal/bayesian"
	"github.com/yourneighborhoodchef/browserforge/internal/data"
)

type HeaderGenerator struct {
	inputNetwork  *bayesian.BayesianNetwork
	headerNetwork *bayesian.BayesianNetwork

	headersOrder map[string][]string

	uniqueBrowsers []string
}

func pascalize(name string) string {

	if strings.HasPrefix(name, ":") || strings.HasPrefix(name, "sec-ch-ua") {
		return name
	}
	lower := strings.ToLower(name)

	switch lower {
	case "dnt", "rtt", "ect":
		return strings.ToUpper(lower)
	}
	parts := strings.Split(name, "-")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = title(p)
	}
	return strings.Join(parts, "-")
}

func title(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func NewHeaderGenerator() (*HeaderGenerator, error) {
	inNet, err := bayesian.LoadInputNetwork()
	if err != nil {
		return nil, fmt.Errorf("loading input network: %w", err)
	}
	hNet, err := bayesian.LoadHeaderNetwork()
	if err != nil {
		return nil, fmt.Errorf("loading header network: %w", err)
	}
	var order map[string][]string
	if err := json.Unmarshal(data.HeadersOrder, &order); err != nil {
		return nil, fmt.Errorf("unmarshalling headers order: %w", err)
	}
	var unique []string
	if err := json.Unmarshal(data.BrowserHelperFile, &unique); err != nil {
		return nil, fmt.Errorf("unmarshalling browser helper file: %w", err)
	}
	return &HeaderGenerator{
		inputNetwork:   inNet,
		headerNetwork:  hNet,
		headersOrder:   order,
		uniqueBrowsers: unique,
	}, nil
}

func (hg *HeaderGenerator) Generate() (map[string]string, error) {

	return hg.GenerateWithConstraints(nil, nil)
}

func (hg *HeaderGenerator) GenerateWithConstraints(
	inputNetConstraints map[string]string,
	requestDependent map[string]string,
) (map[string]string, error) {

	inSample := make(map[string]string)
	if inputNetConstraints != nil {
		for k, v := range inputNetConstraints {
			inSample[k] = v
		}
	}
	inputSample, err := hg.inputNetwork.GenerateSample(inSample)
	if err != nil {
		return nil, fmt.Errorf("sampling input network: %w", err)
	}

	if requestDependent != nil {
		for k, v := range requestDependent {

			inputSample[k] = v
		}
	}

	sample, err := hg.headerNetwork.GenerateSample(inputSample)
	if err != nil {
		return nil, fmt.Errorf("sampling header network: %w", err)
	}

	const missingToken = "*MISSING_VALUE*"
	filtered := make(map[string]string, len(sample))
	for key, val := range sample {
		if strings.HasPrefix(key, "*") {
			continue
		}
		if strings.EqualFold(key, "connection") && val == "close" {
			continue
		}
		if val == missingToken {
			continue
		}
		filtered[key] = val
	}
	result := make(map[string]string, len(filtered))
	for key, val := range filtered {
		result[pascalize(key)] = val
	}
	return result, nil
}
