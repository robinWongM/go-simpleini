package simpleini

import "testing"

func TestParseFromString(t *testing.T) {
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

	configuration, err := parseFromString(testIni)

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
