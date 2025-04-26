// Package browserforge provides browser fingerprinting capabilities
// for generating realistic browser profiles.
package browserforge

import (
	"github.com/yourneighborhoodchef/browserforge/fingerprint"
)

// Generator is the main fingerprint generator
type Generator = fingerprint.Generator

// Fingerprint represents a complete browser fingerprint
type Fingerprint = fingerprint.Fingerprint

// ScreenFingerprint represents screen-related fingerprint data
type ScreenFingerprint = fingerprint.ScreenFingerprint

// NavigatorFingerprint represents navigator-related fingerprint data
type NavigatorFingerprint = fingerprint.NavigatorFingerprint

// VideoCard represents video card information in a fingerprint
type VideoCard = fingerprint.VideoCard

// Option is a function that modifies a Generator
type Option = fingerprint.Option

// New creates a new fingerprint generator with default settings
func New() (*Generator, error) {
	return fingerprint.New()
}

// NewWithOptions creates a new Generator with the specified options
func NewWithOptions(opts ...Option) (*Generator, error) {
	return fingerprint.NewWithOptions(opts...)
}

// WithCustomUserAgent allows specifying a user agent manually
func WithCustomUserAgent(userAgent string) Option {
	return fingerprint.WithCustomUserAgent(userAgent)
}

// WithSeed sets a specific random seed for deterministic fingerprint generation
func WithSeed(seed int64) Option {
	return fingerprint.WithSeed(seed)
}

// WithDeviceCategory allows selecting a device category (desktop, mobile, tablet)
func WithDeviceCategory(category string) Option {
	return fingerprint.WithDeviceCategory(category)
}

// WithBrowser allows selecting a specific browser to emulate
func WithBrowser(browser string) Option {
	return fingerprint.WithBrowser(browser)
}

// WithOperatingSystem allows selecting a specific OS to emulate
func WithOperatingSystem(os string) Option {
	return fingerprint.WithOperatingSystem(os)
}

// WithCamoufoxConstraints sets all Camoufox-compatible constraints
func WithCamoufoxConstraints() Option {
	return fingerprint.WithCamoufoxConstraints()
}

// WithScreenConstraints sets max width and height constraints for screen dimensions
func WithScreenConstraints(maxWidth, maxHeight int) Option {
	return fingerprint.WithScreenConstraints(maxWidth, maxHeight)
}

// WithWindowSize sets a specific outer window size
func WithWindowSize(width, height int) Option {
	return fingerprint.WithWindowSize(width, height)
}