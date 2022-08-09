package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDataUserAgent = []struct {
	input, expected string
	isMobile, isApp bool
}{
	{
		input:    "BhinnekaApp/1.0.0-alpha (Android 9; ASUS_X00TD)",
		expected: "BhinnekaApp/1.0.0-alpha (Android 9; ASUS_X00TD)",
		isMobile: true,
		isApp:    true,
	},
	{
		input:    "BhinnekaApp/1.0.0-alpha (iOS14; iphone 12)",
		expected: "BhinnekaApp/1.0.0-alpha (iOS14; iphone 12)",
		isMobile: true,
		isApp:    true,
	},
	{
		input:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36",
		expected: "Chrome 89.0.4389.82 - macOS 11.1.0",
		isMobile: false,
		isApp:    false,
	},
	{
		input:    "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Mobile Safari/537.36",
		expected: "Chrome 85.0.4183.83 - Android 6.0 [Nexus 5]",
		isMobile: true,
		isApp:    false,
	},
	{
		input:    "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1",
		expected: "Safari 10.0 - iOS 10.3.1 [iPhone]",
		isMobile: true,
		isApp:    false,
	},
	{
		input:    "Mozilla/5.0 (Linux; Android 9; SM-A507FN) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Mobile Safari/537.36",
		expected: "Chrome 81.0.4044.138 - Android 9",
		isMobile: true,
		isApp:    false,
	},
	{
		input:    "Mozilla/5.0 (Linux; U; Android 9; en-us; POCOPHONE F1 Build/PKQ1.180729.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/71.0.3578.141 Mobile Safari/537.36 XiaoMi/MiuiBrowser/12.5.2-go",
		expected: "Chrome 71.0.3578.141 - Android 9 [POCOPHONE F1]",
		isMobile: true,
		isApp:    false,
	},
	{
		input:    "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36",
		expected: "Chrome 78.0.3904.87 - Windows 6.1",
		isMobile: false,
		isApp:    false,
	},
	{
		input:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36",
		expected: "Chrome 80.0.3987.87 - Windows 10.0",
		isMobile: false,
		isApp:    false,
	},
	{
		input:    "Mozilla/5.0 (Windows NT 6.1; ) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
		expected: "Chrome 83.0.4103.116 - Windows 6.1",
		isMobile: false,
	},
	{
		input:    "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0",
		expected: "Firefox 82.0 - Windows 6.1",
		isMobile: false,
		isApp:    false,
	},
	{
		input:    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/80.0.3987.87 HeadlessChrome/80.0.3987.87 Safari/537.36",
		expected: "Safari 537.36 - Linux x86_64",
		isMobile: false,
		isApp:    false,
	},
	{
		input:    "PostmanRuntime/7.17.1",
		expected: "PostmanRuntime 7.17.1",
		isMobile: false,
		isApp:    false,
	},
}

func TestParseUserAgent(t *testing.T) {
	for _, tc := range testDataUserAgent {
		output, isMobile, isApp := ParseUserAgent(tc.input)
		assert.Equal(t, tc.expected, output)
		assert.Equal(t, tc.isMobile, isMobile)
		assert.Equal(t, tc.isApp, isApp)
	}
}
