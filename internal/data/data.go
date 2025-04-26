package data

import (
	"embed"
)

// All embedded data files
//
//go:embed *.json
var files embed.FS

// InputNetwork contains the raw JSON for the input Bayesian network.
var InputNetwork []byte

// HeaderNetwork contains the raw JSON for the header Bayesian network.
var HeaderNetwork []byte

// FingerprintNetwork contains the raw JSON for the fingerprint Bayesian network.
var FingerprintNetwork []byte

// HeadersOrder contains the raw JSON mapping browsers to their expected header order.
var HeadersOrder []byte

// BrowserHelperFile contains the raw JSON with browser information.
var BrowserHelperFile []byte

// init loads all embedded files into variables
func init() {
	var err error
	
	// Load input network
	InputNetwork, err = files.ReadFile("input-network.json")
	if err != nil {
		panic("failed to load input-network.json: " + err.Error())
	}
	
	// Load header network
	HeaderNetwork, err = files.ReadFile("header-network.json")
	if err != nil {
		panic("failed to load header-network.json: " + err.Error())
	}
	
	// Load fingerprint network
	FingerprintNetwork, err = files.ReadFile("fingerprint-network.json")
	if err != nil {
		panic("failed to load fingerprint-network.json: " + err.Error())
	}
	
	// Load headers order
	HeadersOrder, err = files.ReadFile("headers-order.json")
	if err != nil {
		panic("failed to load headers-order.json: " + err.Error())
	}
	
	// Load browser helper file
	BrowserHelperFile, err = files.ReadFile("browser-helper-file.json")
	if err != nil {
		panic("failed to load browser-helper-file.json: " + err.Error())
	}
}