package simpleini

import (
	"reflect"
	"testing"
)

func TestParseFromStringWithCommentDelimiter(t *testing.T) {
	testIni := `# possible values : production, development
app_mode = development

[paths]
# Path to where grafana can store temp files, sessions, and the sqlite3 db (if that is used)
data = /home/git/grafana

[server]
# Protocol (http or https)
protocol = http    # It's recommended to use https

# The http port  to use
http_port = 9999

# Redirect to correct domain if host header does not match domain
# Prevents DNS rebinding attacks
enforce_domain = true`

	testCases := []struct {
		Section string
		Key     string
		Value   string
	}{
		{"", "app_mode", "development"},
		{"paths", "data", "/home/git/grafana"},
		{"server", "protocol", "http"},
		{"server", "http_port", "9999"},
		{"server", "enforce_domain", "true"},
	}

	configuration, err := parseFromStringWithCommentDelimiter(testIni, "#")

	if err != nil {
		t.Errorf("ParseFromString failed, got error %e", err)
	}

	for _, testCase := range testCases {
		got := configuration.Get(testCase.Section, testCase.Key)
		expected := testCase.Value
		if expected != got {
			t.Errorf("ParseFromString wrong, expected %v, got %v", expected, got)
		}
	}
}

func TestParseLineUnix(t *testing.T) {
	testParseLine(t, "#")
}

func TestParseLineWindows(t *testing.T) {
	testParseLine(t, ";")
}

func testParseLine(t *testing.T, commentDelimiter string) {
	testCases := []struct {
		name     string
		iniLine  string
		expected configurationLine
	}{
		{"key value pair line", "app_mode = development", configurationLine{lineKeyValue, "app_mode", "development"}},
		{"key value pair line without spaces", "app_mode=production", configurationLine{lineKeyValue, "app_mode", "production"}},
		{"value contains =", "env=NODE_ENV=production", configurationLine{lineKeyValue, "env", "NODE_ENV=production"}},
		{"comment", "; test", configurationLine{lineBlank, "", ""}},
		{"comment contains key pair line", commentDelimiter + "app_mode = development", configurationLine{lineBlank, "", ""}},
		{"key pair line with comment", "app_mode = development " + commentDelimiter + "set app_mode", configurationLine{lineKeyValue, "app_mode", "development"}},
		{"section", "[section]", configurationLine{lineSection, "", "section"}},
		{"comment contains section", commentDelimiter + "[section]", configurationLine{lineBlank, "", ""}},
		{"section with spaces", " [ test ]  ", configurationLine{lineSection, "", "test"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			parsedLine := parseLine(testCase.iniLine, commentDelimiter)

			if !reflect.DeepEqual(testCase.expected, parsedLine) {
				t.Errorf("expected %v, got %v", testCase.expected, parsedLine)
			}
		})
	}
}
