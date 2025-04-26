package fingerprint

import (
	"fmt"
	"math/rand"
)

// Option is a function that modifies a Generator.
type Option func(*Generator) error

// WithCustomUserAgent allows specifying a user agent manually.
// The rest of the fingerprint will be generated to match this user agent.
func WithCustomUserAgent(userAgent string) Option {
	return func(g *Generator) error {
		g.customUserAgent = userAgent
		return nil
	}
}

// WithSeed sets a specific random seed for deterministic fingerprint generation.
func WithSeed(seed int64) Option {
	return func(g *Generator) error {
		g.seed = &seed
		return nil
	}
}

// WithDeviceCategory allows selecting a device category (desktop, mobile, tablet).
func WithDeviceCategory(category string) Option {
	return func(g *Generator) error {
		g.deviceOption = category
		return nil
	}
}

// WithBrowser allows selecting a specific browser to emulate.
func WithBrowser(browser string) Option {
	return func(g *Generator) error {
		g.browserOption = browser
		return nil
	}
}

// WithOperatingSystem allows selecting a specific OS to emulate.
func WithOperatingSystem(os string) Option {
	return func(g *Generator) error {
		g.osOption = os
		return nil
	}
}

// WithCamoufoxConstraints sets all Camoufox-compatible constraints:
// - Firefox browser only
// - Desktop OS only (linux, macos, windows)
// - Whitelist mode for properties
// This follows the constraints in BROWSERFORGE_CONSTRAINTS.md
func WithCamoufoxConstraints() Option {
	return func(g *Generator) error {
		// Fixed browser: firefox only
		g.browserOption = "firefox"
		
		// Fixed OS category: desktop only with random selection
		desktopOS := []string{"linux", "macos", "windows"}
		if g.seed != nil {
			// Use the existing seed for deterministic selection
			r := rand.New(rand.NewSource(*g.seed))
			g.osOption = desktopOS[r.Intn(len(desktopOS))]
		} else {
			// Random selection
			g.osOption = desktopOS[rand.Intn(len(desktopOS))]
		}
		
		// Enable whitelist filtering
		g.enableWhitelist = true
		
		return nil
	}
}

// WithScreenConstraints sets max width and height constraints for screen dimensions.
// This ensures realistic screen dimensions that don't exceed the specified limits.
func WithScreenConstraints(maxWidth, maxHeight int) Option {
	return func(g *Generator) error {
		if maxWidth <= 0 || maxHeight <= 0 {
			return fmt.Errorf("invalid screen constraints: dimensions must be positive")
		}
		g.screenConstraints = &ScreenConstraints{
			MaxWidth:  maxWidth,
			MaxHeight: maxHeight,
		}
		return nil
	}
}

// WithWindowSize sets a specific outer window size.
// This will override the generated window size in the fingerprint.
func WithWindowSize(width, height int) Option {
	return func(g *Generator) error {
		if width <= 0 || height <= 0 {
			return fmt.Errorf("invalid window size: dimensions must be positive")
		}
		g.windowSize = &WindowSize{
			Width:  width,
			Height: height,
		}
		return nil
	}
}

// NewWithOptions creates a new Generator with the specified options.
func NewWithOptions(opts ...Option) (*Generator, error) {
	// Create the default generator
	g, err := New()
	if err != nil {
		return nil, err
	}
	
	// Apply each option
	for _, opt := range opts {
		if err := opt(g); err != nil {
			return nil, err
		}
	}
	
	return g, nil
}