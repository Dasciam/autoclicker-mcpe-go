package application

import (
	"regexp"
	"strings"
)

// WindowMatcher represents a matcher for the name of the windows.
type WindowMatcher interface {
	// Match returns whether the window name is
	// acceptable for auto-clicker
	Match(title string) bool
}

type WindowMatchStrategy byte

func (w WindowMatchStrategy) String() string {
	switch w {
	case MatchStrategyRegex:
		return "RegExp"
	case MatchStrategyExact:
		return "Exact string"
	case MatchStrategyContains:
		return "Contains string"
	default:
		panic("should never happen")
	}
}

func (w WindowMatchStrategy) New(str string) (WindowMatcher, error) {
	switch w {
	case MatchStrategyRegex:
		regex, err := regexp.Compile(str)
		if err != nil {
			return nil, err
		}
		return &regexStringMatcher{regex: regex}, nil
	case MatchStrategyExact:
		return &exactStringMatcher{match: str}, nil
	case MatchStrategyContains:
		return &containsStringMatcher{match: str}, nil
	default:
		panic("should never happen")
	}
}

const (
	MatchStrategyExact WindowMatchStrategy = iota
	MatchStrategyRegex
	MatchStrategyContains
)

var windowMatchStrategies = [...]WindowMatchStrategy{
	0: MatchStrategyExact,
	1: MatchStrategyRegex,
	2: MatchStrategyContains,
}

type exactStringMatcher struct {
	match string
}

func (e *exactStringMatcher) Match(title string) bool {
	return e.match == title
}

type containsStringMatcher struct {
	match string
}

func (c *containsStringMatcher) Match(title string) bool {
	return strings.Contains(title, c.match)
}

type regexStringMatcher struct {
	regex *regexp.Regexp
}

func (r *regexStringMatcher) Match(title string) bool {
	return r.regex.MatchString(title)
}
