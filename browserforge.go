package browserforge

import (
	"github.com/yourneighborhoodchef/browserforge/fingerprint"
)

type Generator = fingerprint.Generator

type Fingerprint = fingerprint.Fingerprint

type ScreenFingerprint = fingerprint.ScreenFingerprint

type NavigatorFingerprint = fingerprint.NavigatorFingerprint

type VideoCard = fingerprint.VideoCard

type Option = fingerprint.Option

func New() (*Generator, error) {
	return fingerprint.New()
}

func NewWithOptions(opts ...Option) (*Generator, error) {
	return fingerprint.NewWithOptions(opts...)
}

func WithCustomUserAgent(userAgent string) Option {
	return fingerprint.WithCustomUserAgent(userAgent)
}

func WithSeed(seed int64) Option {
	return fingerprint.WithSeed(seed)
}

func WithDeviceCategory(category string) Option {
	return fingerprint.WithDeviceCategory(category)
}

func WithBrowser(browser string) Option {
	return fingerprint.WithBrowser(browser)
}

func WithOperatingSystem(os string) Option {
	return fingerprint.WithOperatingSystem(os)
}

func WithCamoufoxConstraints() Option {
	return fingerprint.WithCamoufoxConstraints()
}

func WithScreenConstraints(maxWidth, maxHeight int) Option {
	return fingerprint.WithScreenConstraints(maxWidth, maxHeight)
}

func WithWindowSize(width, height int) Option {
	return fingerprint.WithWindowSize(width, height)
}
