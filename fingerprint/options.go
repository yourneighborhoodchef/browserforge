package fingerprint

import (
	"fmt"
	"math/rand"
)

type Option func(*Generator) error

func WithCustomUserAgent(userAgent string) Option {
	return func(g *Generator) error {
		g.customUserAgent = userAgent
		return nil
	}
}

func WithSeed(seed int64) Option {
	return func(g *Generator) error {
		g.seed = &seed
		return nil
	}
}

func WithDeviceCategory(category string) Option {
	return func(g *Generator) error {
		g.deviceOption = category
		return nil
	}
}

func WithBrowser(browser string) Option {
	return func(g *Generator) error {
		g.browserOption = browser
		return nil
	}
}

func WithOperatingSystem(os string) Option {
	return func(g *Generator) error {
		g.osOption = os
		return nil
	}
}

func WithCamoufoxConstraints() Option {
	return func(g *Generator) error {

		g.browserOption = "firefox"

		desktopOS := []string{"linux", "macos", "windows"}
		if g.seed != nil {

			r := rand.New(rand.NewSource(*g.seed))
			g.osOption = desktopOS[r.Intn(len(desktopOS))]
		} else {

			g.osOption = desktopOS[rand.Intn(len(desktopOS))]
		}

		g.enableWhitelist = true

		return nil
	}
}

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

func NewWithOptions(opts ...Option) (*Generator, error) {

	g, err := New()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(g); err != nil {
			return nil, err
		}
	}

	return g, nil
}
