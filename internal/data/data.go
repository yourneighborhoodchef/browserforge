package data

import (
	"embed"
)

var files embed.FS

var InputNetwork []byte

var HeaderNetwork []byte

var FingerprintNetwork []byte

var HeadersOrder []byte

var BrowserHelperFile []byte

func init() {
	var err error

	InputNetwork, err = files.ReadFile("input-network.json")
	if err != nil {
		panic("failed to load input-network.json: " + err.Error())
	}

	HeaderNetwork, err = files.ReadFile("header-network.json")
	if err != nil {
		panic("failed to load header-network.json: " + err.Error())
	}

	FingerprintNetwork, err = files.ReadFile("fingerprint-network.json")
	if err != nil {
		panic("failed to load fingerprint-network.json: " + err.Error())
	}

	HeadersOrder, err = files.ReadFile("headers-order.json")
	if err != nil {
		panic("failed to load headers-order.json: " + err.Error())
	}

	BrowserHelperFile, err = files.ReadFile("browser-helper-file.json")
	if err != nil {
		panic("failed to load browser-helper-file.json: " + err.Error())
	}
}
