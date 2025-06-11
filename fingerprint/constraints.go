package fingerprint

import (
	"math/rand"
	"regexp"
)

func applyScreenConstraints(screen *ScreenFingerprint, constraints *ScreenConstraints) {
	if constraints == nil {
		return
	}

	if screen.Width > constraints.MaxWidth {
		screen.Width = constraints.MaxWidth
	}
	if screen.Height > constraints.MaxHeight {
		screen.Height = constraints.MaxHeight
	}

	if screen.AvailWidth > constraints.MaxWidth {
		screen.AvailWidth = constraints.MaxWidth
	}
	if screen.AvailHeight > constraints.MaxHeight {
		screen.AvailHeight = constraints.MaxHeight
	}

	if screen.InnerWidth > screen.Width {
		screen.InnerWidth = screen.Width
	}
	if screen.InnerHeight > screen.Height {
		screen.InnerHeight = screen.Height
	}
}

func applyWindowSize(screen *ScreenFingerprint, windowSize *WindowSize) {
	if windowSize == nil {
		return
	}

	screen.ScreenX += (screen.Width - windowSize.Width) / 2
	screen.PageXOffset = (screen.Width - windowSize.Width) / 2

	if screen.InnerWidth > 0 {
		screen.InnerWidth = max(windowSize.Width-screen.OuterWidth+screen.InnerWidth, 0)
	}
	if screen.InnerHeight > 0 {
		screen.InnerHeight = max(windowSize.Height-screen.OuterHeight+screen.InnerHeight, 0)
	}

	screen.OuterWidth = windowSize.Width
	screen.OuterHeight = windowSize.Height
}

func handleScreenPositioning(screen *ScreenFingerprint) {
	sx := screen.ScreenX

	if sx == 0 {
		screen.ScreenX = 0
		screen.PageXOffset = 0
		screen.PageYOffset = 0
		return
	}

	if sx >= -50 && sx <= 50 {
		screen.PageYOffset = sx
		return
	}

	maxY := screen.AvailHeight - screen.OuterHeight

	if maxY == 0 {
		screen.PageYOffset = 0
	} else if maxY > 0 {
		screen.PageYOffset = rand.Intn(maxY)
	} else {

		screen.PageYOffset = maxY + rand.Intn(-maxY)
	}
}

func filterFalsyValues(fp *Fingerprint) {

	if fp.Screen.AvailHeight < 0 {
		fp.Screen.AvailHeight = 0
	}
	if fp.Screen.AvailWidth < 0 {
		fp.Screen.AvailWidth = 0
	}
	if fp.Screen.AvailTop < 0 {
		fp.Screen.AvailTop = 0
	}
	if fp.Screen.AvailLeft < 0 {
		fp.Screen.AvailLeft = 0
	}
	if fp.Screen.Height < 0 {
		fp.Screen.Height = 0
	}
	if fp.Screen.Width < 0 {
		fp.Screen.Width = 0
	}
	if fp.Screen.InnerHeight < 0 {
		fp.Screen.InnerHeight = 0
	}
	if fp.Screen.InnerWidth < 0 {
		fp.Screen.InnerWidth = 0
	}
	if fp.Screen.OuterHeight < 0 {
		fp.Screen.OuterHeight = 0
	}
	if fp.Screen.OuterWidth < 0 {
		fp.Screen.OuterWidth = 0
	}
	if fp.Screen.ScreenX < 0 {
		fp.Screen.ScreenX = 0
	}
	if fp.Screen.PageXOffset < 0 {
		fp.Screen.PageXOffset = 0
	}
	if fp.Screen.PageYOffset < 0 {
		fp.Screen.PageYOffset = 0
	}
}

func updateFirefoxVersion(fp *Fingerprint, realVersion string) {

	if realVersion == "" {
		return
	}

	re := regexp.MustCompile(`(?<!\d)(1[0-9]{2})(\.[0-9]+)(?!\d)`)

	fp.Navigator.UserAgent = re.ReplaceAllString(fp.Navigator.UserAgent, realVersion+"$2")

	fp.Navigator.AppVersion = re.ReplaceAllString(fp.Navigator.AppVersion, realVersion+"$2")

	if fp.Navigator.Oscpu != nil {
		*fp.Navigator.Oscpu = re.ReplaceAllString(*fp.Navigator.Oscpu, realVersion+"$2")
	}
}

func whitelistProperties(fp *Fingerprint, whitelist PropertyWhitelist) *Fingerprint {

	result := &Fingerprint{
		Headers: make(map[string]string),
		Battery: make(map[string]interface{}),
	}

	screenData := ScreenFingerprint{}
	for _, prop := range whitelist.Screen {
		switch prop {
		case "availHeight":
			screenData.AvailHeight = fp.Screen.AvailHeight
		case "availWidth":
			screenData.AvailWidth = fp.Screen.AvailWidth
		case "availTop":
			screenData.AvailTop = fp.Screen.AvailTop
		case "availLeft":
			screenData.AvailLeft = fp.Screen.AvailLeft
		case "width":
			screenData.Width = fp.Screen.Width
		case "height":
			screenData.Height = fp.Screen.Height
		case "colorDepth":
			screenData.ColorDepth = fp.Screen.ColorDepth
		case "pixelDepth":
			screenData.PixelDepth = fp.Screen.PixelDepth
		case "pageXOffset":
			screenData.PageXOffset = fp.Screen.PageXOffset
		case "pageYOffset":
			screenData.PageYOffset = fp.Screen.PageYOffset
		case "outerWidth":
			screenData.OuterWidth = fp.Screen.OuterWidth
		case "outerHeight":
			screenData.OuterHeight = fp.Screen.OuterHeight
		case "innerWidth":
			screenData.InnerWidth = fp.Screen.InnerWidth
		case "innerHeight":
			screenData.InnerHeight = fp.Screen.InnerHeight
		case "screenX":
			screenData.ScreenX = fp.Screen.ScreenX
		case "screenY":

			screenData.PageYOffset = fp.Screen.PageYOffset
		}
	}
	result.Screen = screenData

	navigatorData := NavigatorFingerprint{}
	for _, prop := range whitelist.Navigator {
		switch prop {
		case "userAgent":
			navigatorData.UserAgent = fp.Navigator.UserAgent
		case "doNotTrack":
			navigatorData.DoNotTrack = fp.Navigator.DoNotTrack
		case "appCodeName":
			navigatorData.AppCodeName = fp.Navigator.AppCodeName
		case "appName":
			navigatorData.AppName = fp.Navigator.AppName
		case "appVersion":
			navigatorData.AppVersion = fp.Navigator.AppVersion
		case "oscpu":
			navigatorData.Oscpu = fp.Navigator.Oscpu
		case "platform":
			navigatorData.Platform = fp.Navigator.Platform
		case "hardwareConcurrency":
			navigatorData.HardwareConcurrency = fp.Navigator.HardwareConcurrency
		case "product":
			navigatorData.Product = fp.Navigator.Product
		case "maxTouchPoints":
			navigatorData.MaxTouchPoints = fp.Navigator.MaxTouchPoints
		case "globalPrivacyControl":
			navigatorData.GlobalPrivacyControl = fp.Navigator.GlobalPrivacyControl
		}
	}
	result.Navigator = navigatorData

	for _, header := range whitelist.Headers {
		if value, exists := fp.Headers[header]; exists {
			result.Headers[header] = value
		}
	}

	if fp.Battery != nil {
		for _, prop := range whitelist.Battery {
			if value, exists := fp.Battery[prop]; exists {
				result.Battery[prop] = value
			}
		}
	}

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func extractFirefoxVersion(userAgent string) string {

	re := regexp.MustCompile(`Firefox/(\d+\.\d+)`)
	matches := re.FindStringSubmatch(userAgent)
	if len(matches) >= 2 {
		return matches[1]
	}

	re = regexp.MustCompile(`Firefox/(\d+)`)
	matches = re.FindStringSubmatch(userAgent)
	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}
