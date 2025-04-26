package headers

import (
   "encoding/json"
   "fmt"
   "strings"

   "github.com/yourneighborhoodchef/browserforge/internal/bayesian"
   "github.com/yourneighborhoodchef/browserforge/internal/data"
)

// HeaderGenerator generates HTTP headers based on constraints.
// HeaderGenerator generates HTTP headers based on the Bayesian networks and helper data.
type HeaderGenerator struct {
   inputNetwork   *bayesian.BayesianNetwork
   headerNetwork  *bayesian.BayesianNetwork
   // headersOrder maps browser names to their ordered list of HTTP header keys.
   headersOrder   map[string][]string
   // uniqueBrowsers holds the raw HTTP browser strings.
   uniqueBrowsers []string
}

// pascalize converts header keys to HTTP standard casing.
func pascalize(name string) string {
   // leave pseudo headers or sec-ch-ua fields intact
   if strings.HasPrefix(name, ":") || strings.HasPrefix(name, "sec-ch-ua") {
       return name
   }
   lower := strings.ToLower(name)
   // uppercase common acronyms
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

// title makes first letter uppercase, rest lowercase.
func title(s string) string {
   if len(s) == 0 {
       return s
   }
   return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// NewHeaderGenerator initializes a new HeaderGenerator.
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

// Generate produces a map of HTTP headers with default constraints.
// It samples from the input and header Bayesian networks, filters out internal fields,
// and pascalizes header names to standard HTTP casing.
func (hg *HeaderGenerator) Generate() (map[string]string, error) {
   // No constraints: empty maps
   return hg.GenerateWithConstraints(nil, nil)
}

// GenerateWithConstraints produces headers given optional input-network constraints
// (e.g., device, operating system, browser HTTP spec) and request-dependent headers
// (e.g., custom User-Agent).
func (hg *HeaderGenerator) GenerateWithConstraints(
   inputNetConstraints map[string]string,
   requestDependent map[string]string,
) (map[string]string, error) {
   // Sample input network with constraints
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
   // Merge in request-dependent headers (e.g., custom User-Agent)
   if requestDependent != nil {
       for k, v := range requestDependent {
           // headerNetwork expects lower-case keys as in network definition
           inputSample[k] = v
       }
   }
   // Sample header network conditioned on inputSample
   sample, err := hg.headerNetwork.GenerateSample(inputSample)
   if err != nil {
       return nil, fmt.Errorf("sampling header network: %w", err)
   }
   // Filter and pascalize
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