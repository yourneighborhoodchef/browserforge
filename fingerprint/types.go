package fingerprint

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

type NavigatorFingerprint struct {
	UserAgent            string                 `json:"userAgent"`
	UserAgentData        map[string]interface{} `json:"userAgentData"`
	DoNotTrack           *string                `json:"doNotTrack,omitempty"`
	AppCodeName          string                 `json:"appCodeName"`
	AppName              string                 `json:"appName"`
	AppVersion           string                 `json:"appVersion"`
	Oscpu                *string                `json:"oscpu,omitempty"`
	Webdriver            bool                   `json:"webdriver"`
	Language             string                 `json:"language"`
	Languages            []string               `json:"languages"`
	Platform             string                 `json:"platform"`
	DeviceMemory         *int                   `json:"deviceMemory,omitempty"`
	HardwareConcurrency  int                    `json:"hardwareConcurrency"`
	Product              string                 `json:"product"`
	ProductSub           string                 `json:"productSub"`
	Vendor               string                 `json:"vendor"`
	VendorSub            *string                `json:"vendorSub,omitempty"`
	MaxTouchPoints       int                    `json:"maxTouchPoints"`
	ExtraProperties      map[string]interface{} `json:"extraProperties"`
	GlobalPrivacyControl *bool                  `json:"globalPrivacyControl,omitempty"`
}

type VideoCard struct {
	Renderer string `json:"renderer"`
	Vendor   string `json:"vendor"`
}

type WindowFingerprint struct {
	InnerHeight      int     `json:"innerHeight"`
	OuterHeight      int     `json:"outerHeight"`
	OuterWidth       int     `json:"outerWidth"`
	InnerWidth       int     `json:"innerWidth"`
	ScreenX          int     `json:"screenX"`
	PageXOffset      int     `json:"pageXOffset"`
	PageYOffset      int     `json:"pageYOffset"`
	DevicePixelRatio float64 `json:"devicePixelRatio"`
}

type WebGLFingerprint struct {
	Renderer string `json:"renderer"`
	Vendor   string `json:"vendor"`
}

type CanvasFingerprint struct {
}

type AudioContextFingerprint struct {
	SampleRate int `json:"sampleRate"`
}

type LocaleFingerprint struct {
	Language  string   `json:"language"`
	Languages []string `json:"languages"`
	TimeZone  string   `json:"timeZone,omitempty"`
}

type Fingerprint struct {
	Screen            ScreenFingerprint      `json:"screen"`
	Navigator         NavigatorFingerprint   `json:"navigator"`
	Headers           map[string]string      `json:"headers"`
	VideoCodecs       map[string]string      `json:"videoCodecs"`
	AudioCodecs       map[string]string      `json:"audioCodecs"`
	PluginsData       map[string]interface{} `json:"pluginsData"`
	Battery           map[string]interface{} `json:"battery,omitempty"`
	VideoCard         *VideoCard             `json:"videoCard,omitempty"`
	MultimediaDevices []string               `json:"multimediaDevices"`
	Fonts             []string               `json:"fonts"`
	MockWebRTC        bool                   `json:"mockWebRTC,omitempty"`
	Slim              bool                   `json:"slim,omitempty"`

	Window       WindowFingerprint       `json:"window"`
	WebGL        WebGLFingerprint        `json:"webgl"`
	Canvas       CanvasFingerprint       `json:"canvas"`
	AudioContext AudioContextFingerprint `json:"audio"`
	Locale       LocaleFingerprint       `json:"locale"`
}

type ScreenConstraints struct {
	MaxWidth  int
	MaxHeight int
}

type WindowSize struct {
	Width  int
	Height int
}

type PropertyWhitelist struct {
	Navigator []string
	Screen    []string
	Headers   []string
	Battery   []string
}

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
