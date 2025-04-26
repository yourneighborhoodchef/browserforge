package fingerprint

// ScreenFingerprint mirrors the Python Screen dataclass
type ScreenFingerprint struct {
	AvailHeight      int     `json:"availHeight"`
	AvailWidth       int     `json:"availWidth"`
	AvailTop         int     `json:"availTop"`
	AvailLeft        int     `json:"availLeft"`
	ColorDepth       int     `json:"colorDepth"`
	Height           int     `json:"height"`
	PixelDepth       int     `json:"pixelDepth"`
	Width            int     `json:"width"`
	DevicePixelRatio float64 `json:"devicePixelRatio"`
	PageXOffset      int     `json:"pageXOffset"`
	PageYOffset      int     `json:"pageYOffset"`
	InnerHeight      int     `json:"innerHeight"`
	OuterHeight      int     `json:"outerHeight"`
	OuterWidth       int     `json:"outerWidth"`
	InnerWidth       int     `json:"innerWidth"`
	ScreenX          int     `json:"screenX"`
	ClientWidth      int     `json:"clientWidth"`
	ClientHeight     int     `json:"clientHeight"`
	HasHDR           bool    `json:"hasHDR"`
}

// NavigatorFingerprint mirrors the Python Navigator dataclass
type NavigatorFingerprint struct {
	UserAgent           string                 `json:"userAgent"`
	UserAgentData       map[string]interface{} `json:"userAgentData"`
	DoNotTrack          *string                `json:"doNotTrack,omitempty"`
	AppCodeName         string                 `json:"appCodeName"`
	AppName             string                 `json:"appName"`
	AppVersion          string                 `json:"appVersion"`
	Oscpu               *string                `json:"oscpu,omitempty"`
	Webdriver           bool                   `json:"webdriver"`
	Language            string                 `json:"language"`
	Languages           []string               `json:"languages"`
	Platform            string                 `json:"platform"`
	DeviceMemory        *int                   `json:"deviceMemory,omitempty"`
	HardwareConcurrency int                    `json:"hardwareConcurrency"`
	Product             string                 `json:"product"`
	ProductSub          string                 `json:"productSub"`
	Vendor              string                 `json:"vendor"`
	VendorSub           *string                `json:"vendorSub,omitempty"`
	MaxTouchPoints      int                    `json:"maxTouchPoints"`
	ExtraProperties     map[string]interface{} `json:"extraProperties"`
	GlobalPrivacyControl *bool                 `json:"globalPrivacyControl,omitempty"`
}

// VideoCard mirrors the Python VideoCard dataclass
type VideoCard struct {
	Renderer string `json:"renderer"`
	Vendor   string `json:"vendor"`
}

// WindowFingerprint defines window-specific properties
type WindowFingerprint struct {
	InnerHeight     int     `json:"innerHeight"`
	OuterHeight     int     `json:"outerHeight"`
	OuterWidth      int     `json:"outerWidth"`
	InnerWidth      int     `json:"innerWidth"`
	ScreenX         int     `json:"screenX"`
	PageXOffset     int     `json:"pageXOffset"`
	PageYOffset     int     `json:"pageYOffset"`
	DevicePixelRatio float64 `json:"devicePixelRatio"`
}

// WebGLFingerprint defines WebGL-specific properties
type WebGLFingerprint struct {
	Renderer string `json:"renderer"`
	Vendor   string `json:"vendor"`
	// Additional WebGL properties can be added here
}

// CanvasFingerprint defines Canvas-specific properties
type CanvasFingerprint struct {
	// Canvas fingerprinting properties
}

// AudioContextFingerprint defines AudioContext-specific properties
type AudioContextFingerprint struct {
	// AudioContext fingerprinting properties
	// Derived from AudioCodecs and other audio-related data
	SampleRate int `json:"sampleRate"`
}

// LocaleFingerprint defines locale-specific properties
type LocaleFingerprint struct {
	Language  string   `json:"language"`
	Languages []string `json:"languages"`
	TimeZone  string   `json:"timeZone,omitempty"`
}

// Fingerprint is the rich result matching the Python version
type Fingerprint struct {
	Screen            ScreenFingerprint       `json:"screen"`
	Navigator         NavigatorFingerprint    `json:"navigator"`
	Headers           map[string]string       `json:"headers"`
	VideoCodecs       map[string]string       `json:"videoCodecs"`
	AudioCodecs       map[string]string       `json:"audioCodecs"`
	PluginsData       map[string]interface{}  `json:"pluginsData"`
	Battery           map[string]interface{}  `json:"battery,omitempty"`
	VideoCard         *VideoCard              `json:"videoCard,omitempty"`
	MultimediaDevices []string                `json:"multimediaDevices"`
	Fonts             []string                `json:"fonts"`
	MockWebRTC        bool                    `json:"mockWebRTC,omitempty"`
	Slim              bool                    `json:"slim,omitempty"`
	
	// Added fields to match the Python interface
	Window         WindowFingerprint     `json:"window"`
	WebGL          WebGLFingerprint      `json:"webgl"`
	Canvas         CanvasFingerprint     `json:"canvas"`
	AudioContext   AudioContextFingerprint `json:"audio"`
	Locale         LocaleFingerprint     `json:"locale"`
}

// ScreenConstraints defines maximum allowed dimensions for screen
type ScreenConstraints struct {
	MaxWidth  int
	MaxHeight int
}

// WindowSize defines a specific window size to override in the fingerprint
type WindowSize struct {
	Width  int
	Height int
}

// PropertyWhitelist defines which properties to include in the fingerprint
// This follows the browserforge.yml whitelist from Camoufox
type PropertyWhitelist struct {
	Navigator []string
	Screen    []string
	Headers   []string
	Battery   []string
}

// DefaultWhitelist returns the Camoufox property whitelist
// This matches the whitelist in BROWSERFORGE_CONSTRAINTS.md
func DefaultWhitelist() PropertyWhitelist {
	return PropertyWhitelist{
		Navigator: []string{
			"userAgent", "doNotTrack", "appCodeName", "appName", "appVersion",
			"oscpu", "platform", "hardwareConcurrency", "product", "maxTouchPoints",
			"globalPrivacyControl",
		},
		Screen: []string{
			"availLeft", "availTop", "availWidth", "availHeight", "width", "height",
			"colorDepth", "pixelDepth", "pageXOffset", "pageYOffset", "outerWidth",
			"outerHeight", "innerWidth", "innerHeight", "screenX", "screenY",
		},
		Headers: []string{
			"Accept-Encoding",
		},
		Battery: []string{
			"charging", "chargingTime", "dischargingTime",
		},
	}
}