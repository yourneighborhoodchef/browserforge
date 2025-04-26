package fingerprint

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/yourneighborhoodchef/browserforge/internal/bayesian"
	"github.com/yourneighborhoodchef/browserforge/internal/headers"
)

// Generator creates browser fingerprints that can be used for web scraping.
type Generator struct {
	network           *bayesian.BayesianNetwork
	headers           *headers.HeaderGenerator
	customUserAgent   string
	seed              *int64
	browserOption     string
	osOption          string
	deviceOption      string
	localeOption      []string
	httpVersionOption string
	strict            bool
	mockWebRTC        bool
	slim              bool
	
	// Camoufox constraints
	enableWhitelist   bool
	screenConstraints *ScreenConstraints
	windowSize        *WindowSize
	firefoxVersion    string
}

// Note: the rich Fingerprint type is defined in types.go

// New creates a new fingerprint generator with default settings.
func New() (*Generator, error) {
	net, err := bayesian.LoadFingerprintNetwork()
	if err != nil {
		return nil, fmt.Errorf("loading fingerprint network: %w", err)
	}
	hg, err := headers.NewHeaderGenerator()
	if err != nil {
		return nil, fmt.Errorf("initializing header generator: %w", err)
	}
	return &Generator{
		network: net,
		headers: hg,
		// by default use non-deterministic seed
		seed: nil,
	}, nil
}

// SetFirefoxVersion allows setting the actual Firefox version to replace placeholder versions
func (g *Generator) SetFirefoxVersion(version string) {
	g.firefoxVersion = version
}

// Generate creates a new fingerprint with matching HTTP headers.
func (g *Generator) Generate() (*Fingerprint, error) {
	// Apply seed for deterministic or random
	if g.seed != nil {
		rand.Seed(*g.seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}
	// Prepare header generation constraints
	inputNet := make(map[string]string)
	if g.browserOption != "" {
		inputNet["*BROWSER"] = g.browserOption
	}
	if g.osOption != "" {
		inputNet["*OPERATING_SYSTEM"] = g.osOption
	}
	if g.deviceOption != "" {
		inputNet["*DEVICE"] = g.deviceOption
	}
	reqDeps := make(map[string]string)
	if g.customUserAgent != "" {
		reqDeps["User-Agent"] = g.customUserAgent
	}
	// Generate headers with constraints
	hdrs, err := g.headers.GenerateWithConstraints(inputNet, reqDeps)
	if err != nil {
		return nil, fmt.Errorf("generating headers: %w", err)
	}
	
	// Extract User-Agent to ensure consistency
	userAgent := hdrs["User-Agent"]
	if userAgent == "" {
		return nil, fmt.Errorf("generated headers missing User-Agent")
	}
	
	// Create constraints with the User-Agent
	constraints := map[string]string{
		"userAgent": userAgent,
	}
	
	// Sample fingerprint network with all constraints and generate rich Fingerprint
	sampleMap, err := g.network.GenerateSample(constraints)
	if err != nil {
		return nil, fmt.Errorf("sampling fingerprint network: %w", err)
	}
	// Transform into rich types
	fp, err := transformFingerprint(sampleMap, hdrs, g.mockWebRTC, g.slim)
	if err != nil {
		return nil, err
	}
	
	// Apply Camoufox constraints if enabled
	if g.enableWhitelist || g.screenConstraints != nil || g.windowSize != nil || g.firefoxVersion != "" {
		fp = g.applyCamoufoxConstraints(fp)
	}
	
	return fp, nil
}

// applyCamoufoxConstraints applies all constraints from BROWSERFORGE_CONSTRAINTS.md
func (g *Generator) applyCamoufoxConstraints(fp *Fingerprint) *Fingerprint {
	// 1. Filter falsy values and clamp negative screen values to 0
	filterFalsyValues(fp)
	
	// 2. Apply screen constraints to limit dimensions to real monitor bounds
	if g.screenConstraints != nil {
		applyScreenConstraints(&fp.Screen, g.screenConstraints)
	}
	
	// 3. Apply window size overrides if specified
	if g.windowSize != nil {
		applyWindowSize(&fp.Screen, g.windowSize)
	}
	
	// 4. Handle screen positioning for realistic window placement
	handleScreenPositioning(&fp.Screen)
	
	// 5. Update Firefox version numbers in various fields
	if g.browserOption == "firefox" && g.firefoxVersion != "" {
		updateFirefoxVersion(fp, g.firefoxVersion)
	} else if fp.Navigator.UserAgent != "" {
		// Try to extract Firefox version from user agent
		detectedVersion := extractFirefoxVersion(fp.Navigator.UserAgent)
		if detectedVersion != "" {
			updateFirefoxVersion(fp, detectedVersion)
		}
	}
	
	// 6. If whitelist is enabled, filter properties to only include those in the whitelist
	if g.enableWhitelist {
		fp = whitelistProperties(fp, DefaultWhitelist())
	}
	
	return fp
}

// GenerateHeadersOnly creates only HTTP headers without the full fingerprint data,
// applying any configured options.
func (g *Generator) GenerateHeadersOnly() (map[string]string, error) {
	// Apply seed
	if g.seed != nil {
		rand.Seed(*g.seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}
	// Prepare constraints
	inputNet := make(map[string]string)
	if g.browserOption != "" {
		inputNet["*BROWSER"] = g.browserOption
	}
	if g.osOption != "" {
		inputNet["*OPERATING_SYSTEM"] = g.osOption
	}
	if g.deviceOption != "" {
		inputNet["*DEVICE"] = g.deviceOption
	}
	reqDeps := make(map[string]string)
	if g.customUserAgent != "" {
		reqDeps["User-Agent"] = g.customUserAgent
	}
	
	headers, err := g.headers.GenerateWithConstraints(inputNet, reqDeps)
	if err != nil {
		return nil, err
	}
	
	// If whitelist enabled, filter headers
	if g.enableWhitelist {
		filteredHeaders := make(map[string]string)
		whitelist := DefaultWhitelist()
		
		// Only keep headers in the whitelist
		for _, allowed := range whitelist.Headers {
			if value, exists := headers[allowed]; exists {
				filteredHeaders[allowed] = value
			}
		}
		
		// Make sure User-Agent is always included
		if _, exists := filteredHeaders["User-Agent"]; !exists {
			if value, exists := headers["User-Agent"]; exists {
				filteredHeaders["User-Agent"] = value
			}
		}
		
		return filteredHeaders, nil
	}
	
	return headers, nil
}

// transformFingerprint maps the raw network sample and headers into the rich Fingerprint type.
func transformFingerprint(
	sample map[string]string,
	headers map[string]string,
	mockWebRTC bool,
	slim bool,
) (*Fingerprint, error) {
	// Create a new fingerprint object
	fp := &Fingerprint{
		Headers:    headers,
		MockWebRTC: mockWebRTC,
		Slim:       slim,
	}

	// Parse the screen data
	if screenStr, ok := sample["screen"]; ok && screenStr != "" {
		if len(screenStr) > len("*STRINGIFIED*") && screenStr[:len("*STRINGIFIED*")] == "*STRINGIFIED*" {
			screenJSON := screenStr[len("*STRINGIFIED*"):]
			var screenData map[string]interface{}
			if err := json.Unmarshal([]byte(screenJSON), &screenData); err != nil {
				return nil, fmt.Errorf("failed to parse screen data: %w", err)
			}
			
			// Map the screen data to the ScreenFingerprint struct
			fp.Screen = ScreenFingerprint{
				AvailHeight:     getIntOrDefault(screenData, "availHeight", 0),
				AvailWidth:      getIntOrDefault(screenData, "availWidth", 0),
				AvailTop:        getIntOrDefault(screenData, "availTop", 0),
				AvailLeft:       getIntOrDefault(screenData, "availLeft", 0),
				ColorDepth:      getIntOrDefault(screenData, "colorDepth", 0),
				Height:          getIntOrDefault(screenData, "height", 0),
				PixelDepth:      getIntOrDefault(screenData, "pixelDepth", 0),
				Width:           getIntOrDefault(screenData, "width", 0),
				DevicePixelRatio: getFloatOrDefault(screenData, "devicePixelRatio", 0),
				PageXOffset:     getIntOrDefault(screenData, "pageXOffset", 0),
				PageYOffset:     getIntOrDefault(screenData, "pageYOffset", 0),
				InnerHeight:     getIntOrDefault(screenData, "innerHeight", 0),
				OuterHeight:     getIntOrDefault(screenData, "outerHeight", 0),
				OuterWidth:      getIntOrDefault(screenData, "outerWidth", 0),
				InnerWidth:      getIntOrDefault(screenData, "innerWidth", 0),
				ScreenX:         getIntOrDefault(screenData, "screenX", 0),
				ClientWidth:     getIntOrDefault(screenData, "clientWidth", 0),
				ClientHeight:    getIntOrDefault(screenData, "clientHeight", 0),
				HasHDR:          getBoolOrDefault(screenData, "hasHDR", false),
			}
		}
	}

	// Parse Navigator data
	userAgent := sample["userAgent"]
	var userAgentData map[string]interface{}
	if uadStr, ok := sample["userAgentData"]; ok && uadStr != "" && uadStr != "*MISSING_VALUE*" {
		if len(uadStr) > len("*STRINGIFIED*") && uadStr[:len("*STRINGIFIED*")] == "*STRINGIFIED*" {
			uadJSON := uadStr[len("*STRINGIFIED*"):]
			if err := json.Unmarshal([]byte(uadJSON), &userAgentData); err != nil {
				return nil, fmt.Errorf("failed to parse userAgentData: %w", err)
			}
		}
	}

	// Parse languages
	var languages []string
	if langStr, ok := sample["languages"]; ok && langStr != "" && langStr != "*MISSING_VALUE*" {
		if len(langStr) > len("*STRINGIFIED*") && langStr[:len("*STRINGIFIED*")] == "*STRINGIFIED*" {
			langJSON := langStr[len("*STRINGIFIED*"):]
			if err := json.Unmarshal([]byte(langJSON), &languages); err != nil {
				return nil, fmt.Errorf("failed to parse languages: %w", err)
			}
		}
	}
	
	// Set language from Accept-Language header if available
	language := ""
	if len(languages) > 0 {
		language = languages[0]
	}
	
	// Get or initialize extraProperties
	var extraProperties map[string]interface{}
	if epStr, ok := sample["extraProperties"]; ok && epStr != "" && epStr != "*MISSING_VALUE*" {
		if len(epStr) > len("*STRINGIFIED*") && epStr[:len("*STRINGIFIED*")] == "*STRINGIFIED*" {
			epJSON := epStr[len("*STRINGIFIED*"):]
			if err := json.Unmarshal([]byte(epJSON), &extraProperties); err != nil {
				return nil, fmt.Errorf("failed to parse extraProperties: %w", err)
			}
		}
	}

	// Construct the NavigatorFingerprint
	navigator := NavigatorFingerprint{
		UserAgent:           userAgent,
		UserAgentData:       userAgentData,
		AppCodeName:         getStringOrDefault(sample, "appCodeName", ""),
		AppName:             getStringOrDefault(sample, "appName", ""),
		AppVersion:          getStringOrDefault(sample, "appVersion", ""),
		Webdriver:           getBoolFromString(sample, "webdriver", false),
		Language:            language,
		Languages:           languages,
		Platform:            getStringOrDefault(sample, "platform", ""),
		HardwareConcurrency: getIntFromString(sample, "hardwareConcurrency", 0),
		Product:             getStringOrDefault(sample, "product", ""),
		ProductSub:          getStringOrDefault(sample, "productSub", ""),
		Vendor:              getStringOrDefault(sample, "vendor", ""),
		MaxTouchPoints:      getIntFromString(sample, "maxTouchPoints", 0),
		ExtraProperties:     extraProperties,
	}
	
	// Handle optional Navigator fields
	if doNotTrack, ok := sample["doNotTrack"]; ok && doNotTrack != "*MISSING_VALUE*" {
		navigator.DoNotTrack = &doNotTrack
	}
	if oscpu, ok := sample["oscpu"]; ok && oscpu != "*MISSING_VALUE*" {
		navigator.Oscpu = &oscpu
	}
	if vendorSub, ok := sample["vendorSub"]; ok && vendorSub != "*MISSING_VALUE*" {
		navigator.VendorSub = &vendorSub
	}
	if deviceMemoryStr, ok := sample["deviceMemory"]; ok && deviceMemoryStr != "*MISSING_VALUE*" {
		deviceMemory, err := strconv.Atoi(deviceMemoryStr)
		if err == nil {
			navigator.DeviceMemory = &deviceMemory
		}
	}
	if gpcStr, ok := sample["globalPrivacyControl"]; ok && gpcStr != "*MISSING_VALUE*" {
		gpc, err := strconv.ParseBool(gpcStr)
		if err == nil {
			navigator.GlobalPrivacyControl = &gpc
		}
	}

	fp.Navigator = navigator

	// Parse VideoCodecs
	var videoCodecs map[string]string
	if vcStr, ok := sample["videoCodecs"]; ok && vcStr != "" && vcStr != "*MISSING_VALUE*" {
		if len(vcStr) > len("*STRINGIFIED*") && vcStr[:len("*STRINGIFIED*")] == "*STRINGIFIED*" {
			vcJSON := vcStr[len("*STRINGIFIED*"):]
			if err := json.Unmarshal([]byte(vcJSON), &videoCodecs); err != nil {
				return nil, fmt.Errorf("failed to parse videoCodecs: %w", err)
			}
		}
	}
	fp.VideoCodecs = videoCodecs

	// Parse AudioCodecs
	var audioCodecs map[string]string
	if acStr, ok := sample["audioCodecs"]; ok && acStr != "" && acStr != "*MISSING_VALUE*" {
		if len(acStr) > len("*STRINGIFIED*") && acStr[:len("*STRINGIFIED*")] == "*STRINGIFIED*" {
			acJSON := acStr[len("*STRINGIFIED*"):]
			if err := json.Unmarshal([]byte(acJSON), &audioCodecs); err != nil {
				return nil, fmt.Errorf("failed to parse audioCodecs: %w", err)
			}
		}
	}
	fp.AudioCodecs = audioCodecs

	// Parse PluginsData
	var pluginsData map[string]interface{}
	if pdStr, ok := sample["pluginsData"]; ok && pdStr != "" && pdStr != "*MISSING_VALUE*" {
		if len(pdStr) > len("*STRINGIFIED*") && pdStr[:len("*STRINGIFIED*")] == "*STRINGIFIED*" {
			pdJSON := pdStr[len("*STRINGIFIED*"):]
			if err := json.Unmarshal([]byte(pdJSON), &pluginsData); err != nil {
				return nil, fmt.Errorf("failed to parse pluginsData: %w", err)
			}
		}
	}
	fp.PluginsData = pluginsData

	// Parse Battery
	var battery map[string]interface{}
	if batteryStr, ok := sample["battery"]; ok && batteryStr != "" && batteryStr != "*MISSING_VALUE*" {
		if len(batteryStr) > len("*STRINGIFIED*") && batteryStr[:len("*STRINGIFIED*")] == "*STRINGIFIED*" {
			batteryJSON := batteryStr[len("*STRINGIFIED*"):]
			if err := json.Unmarshal([]byte(batteryJSON), &battery); err != nil {
				return nil, fmt.Errorf("failed to parse battery: %w", err)
			}
			fp.Battery = battery
		}
	}

	// Parse VideoCard
	if videoCardStr, ok := sample["videoCard"]; ok && videoCardStr != "" && videoCardStr != "*MISSING_VALUE*" {
		if len(videoCardStr) > len("*STRINGIFIED*") && videoCardStr[:len("*STRINGIFIED*")] == "*STRINGIFIED*" {
			videoCardJSON := videoCardStr[len("*STRINGIFIED*"):]
			var videoCardData map[string]interface{}
			if err := json.Unmarshal([]byte(videoCardJSON), &videoCardData); err != nil {
				return nil, fmt.Errorf("failed to parse videoCard: %w", err)
			}
			
			fp.VideoCard = &VideoCard{
				Renderer: fmt.Sprintf("%v", videoCardData["renderer"]),
				Vendor:   fmt.Sprintf("%v", videoCardData["vendor"]),
			}
		}
	}

	// Parse MultimediaDevices
	var multimediaDevices []string
	if mdStr, ok := sample["multimediaDevices"]; ok && mdStr != "" && mdStr != "*MISSING_VALUE*" {
		if len(mdStr) > len("*STRINGIFIED*") && mdStr[:len("*STRINGIFIED*")] == "*STRINGIFIED*" {
			mdJSON := mdStr[len("*STRINGIFIED*"):]
			// First try to parse as string array
			err := json.Unmarshal([]byte(mdJSON), &multimediaDevices)
			if err != nil {
				// If parsing as string array fails, try to parse as map/object and extract values
				var mdMap map[string]interface{}
				if err = json.Unmarshal([]byte(mdJSON), &mdMap); err == nil {
					// Convert map to slice of strings
					for _, v := range mdMap {
						if strValue, ok := v.(string); ok {
							multimediaDevices = append(multimediaDevices, strValue)
						}
					}
				} else {
					return nil, fmt.Errorf("failed to parse multimediaDevices: %w", err)
				}
			}
		}
	}
	fp.MultimediaDevices = multimediaDevices

	// Parse Fonts
	var fonts []string
	if fontsStr, ok := sample["fonts"]; ok && fontsStr != "" && fontsStr != "*MISSING_VALUE*" {
		if len(fontsStr) > len("*STRINGIFIED*") && fontsStr[:len("*STRINGIFIED*")] == "*STRINGIFIED*" {
			fontsJSON := fontsStr[len("*STRINGIFIED*"):]
			if err := json.Unmarshal([]byte(fontsJSON), &fonts); err != nil {
				return nil, fmt.Errorf("failed to parse fonts: %w", err)
			}
		}
	}
	fp.Fonts = fonts

	return fp, nil
}

// Helper functions for type conversion

func getStringOrDefault(m map[string]string, key, defaultVal string) string {
	if val, ok := m[key]; ok && val != "*MISSING_VALUE*" {
		return val
	}
	return defaultVal
}

func getIntFromString(m map[string]string, key string, defaultVal int) int {
	if val, ok := m[key]; ok && val != "*MISSING_VALUE*" {
		intVal, err := strconv.Atoi(val)
		if err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getBoolFromString(m map[string]string, key string, defaultVal bool) bool {
	if val, ok := m[key]; ok && val != "*MISSING_VALUE*" {
		boolVal, err := strconv.ParseBool(val)
		if err == nil {
			return boolVal
		}
	}
	return defaultVal
}

func getIntOrDefault(m map[string]interface{}, key string, defaultVal int) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case string:
			intVal, err := strconv.Atoi(v)
			if err == nil {
				return intVal
			}
		}
	}
	return defaultVal
}

func getFloatOrDefault(m map[string]interface{}, key string, defaultVal float64) float64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case string:
			floatVal, err := strconv.ParseFloat(v, 64)
			if err == nil {
				return floatVal
			}
		}
	}
	return defaultVal
}

func getBoolOrDefault(m map[string]interface{}, key string, defaultVal bool) bool {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case bool:
			return v
		case string:
			boolVal, err := strconv.ParseBool(v)
			if err == nil {
				return boolVal
			}
		}
	}
	return defaultVal
}