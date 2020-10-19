package simpleini

import (
	"regexp"
	"runtime"
	"strings"
)

type ConfigurationSection map[string]string

type Configuration map[string]ConfigurationSection

func (c *Configuration) Get(section, key string) string {
	return (*c)[section][key]
}

type LineType string

type ConfigurationLine struct {
	LineType LineType
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
	sections := make(map[string]ConfigurationSection)

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

func ParseLine(iniLine, commentDelimiter string) ConfigurationLine {
	// remove comments
	iniLine = strings.SplitN(iniLine, commentDelimiter, 2)[0]

	sectionRegExp := regexp.MustCompile("^\\s*\\[\\s*(.*?)\\s*\\]\\s*$")
	keyValueRegExp := regexp.MustCompile("^\\s*(.*?)\\s*=\\s*(.*?)\\s*$")

	sectionFound := sectionRegExp.FindStringSubmatch(iniLine)

	if len(sectionFound) == 2 {
		return ConfigurationLine{lineSection, "", sectionFound[1]}
	}

	keyValueFound := keyValueRegExp.FindStringSubmatch(iniLine)

	if len(keyValueFound) == 3 {
		return ConfigurationLine{lineKeyValue, keyValueFound[1], keyValueFound[2]}
	}

	return ConfigurationLine{lineBlank, "", ""}
}
