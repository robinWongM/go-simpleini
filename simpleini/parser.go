package simpleini

import (
	"regexp"
	"runtime"
	"strings"
)

type configurationSection map[string]string

// Configuration provides a method `GET` to access configuration fields.
type Configuration map[string]configurationSection

// Get a configuration field by providing section name and key name.
// The default section is "".
func (c *Configuration) Get(section, key string) string {
	return (*c)[section][key]
}

type lineType string

type configurationLine struct {
	LineType lineType
	Key      string
	Value    string
}

const (
	lineBlank    = "lineBlank"
	lineSection  = "lineSection"
	lineKeyValue = "lineKeyValue"
)

var defaultCommentDelimiter string

func init() {
	switch runtime.GOOS {
	case "linux":
		defaultCommentDelimiter = "#"
	case "windows":
		defaultCommentDelimiter = ";"
	}
}

func ParseFromString(iniContent string) (Configuration, error) {
	return ParseFromStringWithCommentDelimiter(iniContent, defaultCommentDelimiter)
}

func ParseFromStringWithCommentDelimiter(iniContent, commentDelimiter string) (Configuration, error) {
	sections := make(map[string]configurationSection)

	currentSection := ""
	sections[currentSection] = make(map[string]string)

	for _, line := range strings.Split(iniContent, "\n") {
		parsedLine := ParseLine(line, commentDelimiter)
		switch parsedLine.LineType {
		case lineSection:
			currentSection = parsedLine.Value
			sections[currentSection] = make(map[string]string)
		case lineKeyValue:
			sections[currentSection][parsedLine.Key] = parsedLine.Value
		}
	}

	return sections, nil
}

func ParseLine(iniLine, commentDelimiter string) configurationLine {
	// remove comments
	iniLine = strings.SplitN(iniLine, commentDelimiter, 2)[0]

	sectionRegExp := regexp.MustCompile("^\\s*\\[\\s*(.*?)\\s*\\]\\s*$")
	keyValueRegExp := regexp.MustCompile("^\\s*(.*?)\\s*=\\s*(.*?)\\s*$")

	sectionFound := sectionRegExp.FindStringSubmatch(iniLine)

	if len(sectionFound) == 2 {
		return configurationLine{lineSection, "", sectionFound[1]}
	}

	keyValueFound := keyValueRegExp.FindStringSubmatch(iniLine)

	if len(keyValueFound) == 3 {
		return configurationLine{lineKeyValue, keyValueFound[1], keyValueFound[2]}
	}

	return configurationLine{lineBlank, "", ""}
}
