package simpleini

import "testing"

func TestParseFromString(t *testing.T) {
	testIni := `[Application]
UseLiveData=1
;coke=zero
pepsi=diet   ;gag
#stackoverflow=splotchy`

	testCases := []struct {
		Section string
		Key     string
		Value   string
	}{
		{"Application", "UseLiveData", "1"},
		{"Application", "pepsi", "diet"},
		{"Application", "#stackoverflow", "splotchy"},
		{"Application", ";coke", ""},
		{"Application", "coke", ""},
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
